package game

import (
	"errors"
    "fmt"
    "log"
	"remi-game/models"
    "github.com/google/uuid"
    "sync"
)

type PersistentStore interface {
    SaveGame(game models.Game) error
}

type GameManager struct {
	Game models.Game
    Mutex sync.Mutex
    Store PersistentStore
}

func NewGame(store PersistentStore) *GameManager {
	return &GameManager{
		Game: models.Game{
			ID:        uuid.New().String(),
			Status:    models.StateWaiting,
			Deck:      NewDeck(),
			Players:   make([]models.Player, 0),
			TableSets: make([][]models.Card, 0),
		},
        Store: store,
	}
}

func (gm *GameManager) save() {
    if gm.Store != nil {
        go gm.Store.SaveGame(gm.Game)
    }
}

func (gm *GameManager) AddPlayer(name string) (*models.Player, error) {
    // Caller holds lock
	if gm.Game.Status != models.StateWaiting {
		return nil, errors.New("game already in progress")
	}
	if len(gm.Game.Players) >= 4 {
		return nil, errors.New("game is full")
	}

    newPlayerID := uuid.New().String()
	player := models.Player{
		ID:          newPlayerID,
		Name:        name,
		SeatIndex:   len(gm.Game.Players),
		IsConnected: true,
		Hand:        make([]models.Card, 0),
        HasTakenFromPile: false,
        PlayedSets:  make([][]models.Card, 0),
	}

    // Assign Master if first player
    if len(gm.Game.Players) == 0 {
        gm.Game.MasterPlayerID = newPlayerID
    }

	gm.Game.Players = append(gm.Game.Players, player)
    gm.save()
	return &player, nil
}

// StartGame deals 8 cards to Master (who starts), 7 to others.
func (gm *GameManager) StartGame(initiatorID string) error {
    gm.Mutex.Lock()
    defer gm.Mutex.Unlock()

    return gm.StartGameUnlocked(initiatorID)
}

// StartGameUnlocked is the lock-free core of StartGame.
// Caller MUST hold gm.Mutex.
func (gm *GameManager) StartGameUnlocked(initiatorID string) error {
    if gm.Game.Status != models.StateWaiting {
        return errors.New("game already running")
    }

    if initiatorID != gm.Game.MasterPlayerID {
        return errors.New("only master can start game")
    }
    
    // Fill empty slots with Bots
    for len(gm.Game.Players) < 4 {
        botName := fmt.Sprintf("Bot %d", len(gm.Game.Players)+1)
        botPlayer := models.Player{
            ID:          uuid.New().String(),
            Name:        botName,
            SeatIndex:   len(gm.Game.Players),
            IsConnected: true,
            Hand:        make([]models.Card, 0),
        }
        gm.Game.Players = append(gm.Game.Players, botPlayer)
    }

	gm.Game.Status = models.StateInProgress
	
	// Shuffle Deck
	gm.Game.Deck = Shuffle(gm.Game.Deck)
	
    // Deal cards: Master gets 8, others 7
	for i := range gm.Game.Players {
        var handSize int
        if gm.Game.Players[i].ID == gm.Game.MasterPlayerID {
            handSize = 8
        } else {
            handSize = 7
        }

		hand := make([]models.Card, handSize)
		copy(hand, gm.Game.Deck[:handSize])
		gm.Game.Players[i].Hand = hand
		gm.Game.Deck = gm.Game.Deck[handSize:]
        gm.Game.Players[i].HasTakenFromPile = false
        gm.Game.Players[i].HasPlayedSet = false
        gm.Game.Players[i].PlayedSets = make([][]models.Card, 0)
        gm.Game.Players[i].Score = 0
	}

    // Pile starts empty
    gm.Game.Pile = make([]models.Card, 0)
    
    // Master starts
    masterIdx := 0
    for i, p := range gm.Game.Players {
        if p.ID == gm.Game.MasterPlayerID {
            masterIdx = i
            break
        }
    }
	gm.Game.CurrentTurnPlayer = masterIdx
    // Master has 8 cards, so they start by Discarding (skip Draw)
    gm.Game.TurnPhase = models.PhasePlay // Logic allows Play/Discard in this phase
    
    gm.save()
	return nil
}

// DrawCard now accepts source: "DECK" or "PILE" and count (1-3 for PILE)
func (gm *GameManager) DrawCard(playerID string, source string, count int) error {
    if gm.Game.Status != models.StateInProgress {
        return errors.New("game not in progress")
    }

    if count < 1 {
        count = 1
    }

    idx := -1
    for i, p := range gm.Game.Players {
        if p.ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not found")
    }
    if idx != gm.Game.CurrentTurnPlayer {
        return errors.New("not your turn")
    }
    if gm.Game.TurnPhase != models.PhaseDraw {
        return errors.New("not draw phase")
    }
    
    var cardsToDraw []models.Card

    if source == "PILE" {
        // Validation for Pile Draw
        if count > 3 {
            return errors.New("cannot draw more than 3 cards from pile")
        }
        if len(gm.Game.Pile) < count {
            return errors.New("not enough cards in pile")
        }
        
        // Take last 'count' cards (copy to avoid slice aliasing)
        startIndex := len(gm.Game.Pile) - count
        cardsToDraw = make([]models.Card, count)
        copy(cardsToDraw, gm.Game.Pile[startIndex:])

        // RESTRICTION: Must form a specific set (Seri) with hand using ALL drawn cards
        // AND player must have a pair in hand (implied by CanFormSetWith check logic)
        if !CanFormSetWithMultiple(gm.Game.Players[idx].Hand, cardsToDraw) {
            return errors.New("cannot take from pile unless all taken cards form a valid set (Seri) with your hand")
        }

        // Remove from Pile
        gm.Game.Pile = gm.Game.Pile[:startIndex]
        gm.Game.Players[idx].HasTakenFromPile = true 

    } else {
        // Default DECK
        if count != 1 {
            return errors.New("can only draw 1 card from deck")
        }

        // Check if deck is empty
        if len(gm.Game.Deck) == 0 {
             // Game Over: Empty Deck
             // Rules: Game ends, calculate scores for everyone (no winner)
             gm.Game.Status = models.StateFinished
             gm.Game.WinnerID = "" // No specific winner
             gm.calculateScores()
             gm.save()
             
             // We return nil because the action "Draw" successfully caused the game state to change (to Finished).
             // The frontend will receive the updated game state (Finished) and show the Scoreboard.
             return nil
        }
        
        cardsToDraw = []models.Card{gm.Game.Deck[0]} // single element, safe copy
        gm.Game.Deck = append([]models.Card{}, gm.Game.Deck[1:]...)
    }

    gm.Game.Players[idx].Hand = append(gm.Game.Players[idx].Hand, cardsToDraw...)
    gm.Game.TurnPhase = models.PhasePlay // Next phase
    
    gm.save()
    return nil
}

// DrawFromPile lets a player pick a specific card from the pile
// by selecting hand cards that together form a valid set.
// The combined set is played to the table immediately.
func (gm *GameManager) DrawFromPile(playerID string, handCardIDs []string, pileCardID string) error {
    if gm.Game.Status != models.StateInProgress {
        return errors.New("game not in progress")
    }

    idx := -1
    for i, p := range gm.Game.Players {
        if p.ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not found")
    }
    if idx != gm.Game.CurrentTurnPlayer {
        return errors.New("not your turn")
    }
    if gm.Game.TurnPhase != models.PhaseDraw {
        return errors.New("not draw phase")
    }
    if len(handCardIDs) < 2 {
        return errors.New("must select at least 2 cards from hand")
    }

    // Find pile card
    pileIdx := -1
    var pileCard models.Card
    for i, c := range gm.Game.Pile {
        if c.ID == pileCardID {
            pileIdx = i
            pileCard = c
            break
        }
    }
    if pileIdx == -1 {
        return errors.New("card not found in pile")
    }

    // Restrict to last 3 cards of pile
    pileLen := len(gm.Game.Pile)
    if pileIdx < pileLen-3 {
        return errors.New("can only pick from the last 3 cards of the pile")
    }

    // Find hand cards
    handMap := make(map[string]models.Card)
    for _, c := range gm.Game.Players[idx].Hand {
        handMap[c.ID] = c
    }

    var handCards []models.Card
    for _, id := range handCardIDs {
        c, ok := handMap[id]
        if !ok {
            return errors.New("card not found in hand: " + id)
        }
        handCards = append(handCards, c)
    }

    // Combine and validate
    combined := make([]models.Card, 0, len(handCards)+1)
    combined = append(combined, handCards...)
    combined = append(combined, pileCard)

    if !IsValidSet(combined) {
        return errors.New("selected cards do not form a valid set")
    }

    // Remove pile card
    gm.Game.Pile = append(gm.Game.Pile[:pileIdx], gm.Game.Pile[pileIdx+1:]...)

    // Remove hand cards
    removeSet := make(map[string]bool)
    for _, id := range handCardIDs {
        removeSet[id] = true
    }
    newHand := make([]models.Card, 0)
    for _, c := range gm.Game.Players[idx].Hand {
        if !removeSet[c.ID] {
            newHand = append(newHand, c)
        }
    }
    gm.Game.Players[idx].Hand = newHand

    // Play the set to table
    gm.Game.TableSets = append(gm.Game.TableSets, combined)
    gm.Game.Players[idx].PlayedSets = append(gm.Game.Players[idx].PlayedSets, combined)
    gm.Game.Players[idx].HasPlayedSet = true
    gm.Game.Players[idx].HasTakenFromPile = true

    // Advance to PLAY phase
    gm.Game.TurnPhase = models.PhasePlay

    gm.save()
    return nil
}

func (gm *GameManager) PlaySet(playerID string, cards []models.Card) error {
    if gm.Game.Status != models.StateInProgress {
        return errors.New("game not in progress")
    }

    idx := -1
    for i, p := range gm.Game.Players {
        if p.ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not found")
    }
    if idx != gm.Game.CurrentTurnPlayer {
        return errors.New("not your turn")
    }
    if gm.Game.TurnPhase != models.PhasePlay {
        return errors.New("must draw first or already in discard phase")
    }

    if !IsValidSet(cards) {
        return errors.New("invalid set")
    }

    toRemove := make(map[string]bool)
    for _, c := range cards {
        toRemove[c.ID] = true
    }

    playerHandMap := make(map[string]models.Card)
    for _, c := range gm.Game.Players[idx].Hand {
        playerHandMap[c.ID] = c
    }

    for id := range toRemove {
        if _, ok := playerHandMap[id]; !ok {
            return errors.New("player does not have these cards")
        }
    }

    newHand := make([]models.Card, 0)
    for _, c := range gm.Game.Players[idx].Hand {
        if !toRemove[c.ID] {
            newHand = append(newHand, c)
        }
    }
    gm.Game.Players[idx].Hand = newHand

    gm.Game.TableSets = append(gm.Game.TableSets, cards) // Keep for global history if needed, or remove?
    gm.Game.Players[idx].PlayedSets = append(gm.Game.Players[idx].PlayedSets, cards)
    gm.Game.Players[idx].HasPlayedSet = true // Mark as played set
    
    gm.save()
    return nil
}

// DiscardCard handles disarding a card to end turn
func (gm *GameManager) DiscardCard(playerID string, cardID string) error {
    if gm.Game.Status != models.StateInProgress {
        return errors.New("game not in progress")
    }

    idx := -1
    for i, p := range gm.Game.Players {
        if p.ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not found")
    }
    if idx != gm.Game.CurrentTurnPlayer {
        return errors.New("not your turn")
    }
    // Phase must be PLAY (or Discard if we had separate phase, but logic uses PLAY for both)
    // Actually, usually you Draw -> Play Sets -> Discard.
    // So logic checks if you have drawn?
    // In current simplified logic: TurnPhase = PLAY means you can play sets OR discard.
    // If we want to enforce "Must Draw First", we check Phase != DRAW.
    if gm.Game.TurnPhase == models.PhaseDraw {
        return errors.New("must draw a card first")
    }

    // Find card in hand
    cardIdx := -1
    for i, c := range gm.Game.Players[idx].Hand {
        if c.ID == cardID {
            cardIdx = i
            break
        }
    }
    if cardIdx == -1 {
        return errors.New("card not found in hand")
    }

    // Remove from hand
    card := gm.Game.Players[idx].Hand[cardIdx]
    gm.Game.Players[idx].Hand = append(gm.Game.Players[idx].Hand[:cardIdx], gm.Game.Players[idx].Hand[cardIdx+1:]...)

    // Add to Pile
    gm.Game.Pile = append(gm.Game.Pile, card)

    // End Turn -> Next Player
    gm.Game.CurrentTurnPlayer = (gm.Game.CurrentTurnPlayer + 1) % 4
    gm.Game.TurnPhase = models.PhaseDraw // Next player starts in Draw phase
    
    // Reset "HasTakenFromPile" for next player? 
    // Usually this flag is per turn.
    // The struct has `HasTakenFromPile` on Player. 
    // We should reset it for the NEW current player.
    gm.Game.Players[gm.Game.CurrentTurnPlayer].HasTakenFromPile = false
    
    gm.save()
    return nil
}

func (gm *GameManager) DeclareWin(playerID string) error {
    log.Printf("[DeclareWin] Called by player: %s", playerID)
    
    if gm.Game.Status != models.StateInProgress {
        return errors.New("game not in progress")
    }

    idx := -1
    for i, p := range gm.Game.Players {
        if p.ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not found")
    }
    if idx != gm.Game.CurrentTurnPlayer {
         return errors.New("not your turn")
    }
    if gm.Game.TurnPhase == models.PhaseDraw {
        return errors.New("must draw a card first before declaring win")
    }

    log.Printf("[DeclareWin] Player %d hand size: %d", idx, len(gm.Game.Players[idx].Hand))
    for i, c := range gm.Game.Players[idx].Hand {
        log.Printf("[DeclareWin]   Card %d: %s %d (ID: %s)", i, c.Suit, c.Rank, c.ID)
    }
    
    // Check if hand is a winning hand
    isWin := IsWinningHand(gm.Game.Players[idx].Hand)
    log.Printf("[DeclareWin] IsWinningHand result: %v", isWin)
    
    if !isWin {
        return errors.New("hand is not a winning hand (cannot be fully formed into sets with 1 discard)")
    }

    gm.Game.Status = models.StateFinished
    gm.Game.WinnerID = playerID
    
    log.Printf("[DeclareWin] Calculating scores...")
    gm.calculateScores()
    
    log.Printf("[DeclareWin] Saving game state...")
    gm.save()
    log.Printf("[DeclareWin] Success! Winner: %s", playerID)
    return nil
}

func (gm *GameManager) calculateScores() {
    // Winner gets 0? User didn't specify winner points.
    // Losers scored based on cards in hand.
    // 2-10: 5, JQK: 10, A: 15, Joker: -250 (Penalty?)
    // Note: Usually points are BAD in Rummy.
    // User said: "Joker: -250 points... A player who ends game... scores -250 (if holding joker)".
    // "All other players are scored...".
    // I will sum up points as penalties.
    
    for i := range gm.Game.Players {
        if gm.Game.Players[i].ID == gm.Game.WinnerID {
            gm.Game.Players[i].Score = 0
            continue
        }
        
        score := 0
        for _, c := range gm.Game.Players[i].Hand {
            if c.Suit == models.Joker {
                score -= 250 // Negative for holding Joker? 
                // Wait. "Joker: -250 points (if not part of a valid combination)".
                // Usually penalties are POSITIVE in Rummy points (bad).
                // Or maybe user wants SCORE (good).
                // "Scores -250 points" usually means BAD.
                // "Number cards: 5 points each". 
                // If I have 3 cards (2,3,4) = 15 points. IS this good or bad?
                // Rummy: Goal is 0. Points are bad.
                // So holding Joker should be HUGE penalty. -250 is... negative penalty?
                // Context: "Scores -250" might mean -250 to total score?
                // Let's assume Points = Bad.
                // So Score += 5.
                // Joker = -250? That reduces score (Good?).
                // Usually holding Joker is BAD -> +25 points or +50 usually.
                // "Joker: -250 points".
                // Maybe the Game Score is "Points you got".
                // Let's just implement exactly: Score += Value.
                // Joker Value = -250.
                // 2-10 Value = 5.
                // So if I have Joker + 2: Score = -250 + 5 = -245.
                // This seems generous if goal is low score.
                // BUT "Joker... scores -250 points".
                // Maybe user implies "Joker COST 250 points"?
                // "Scores -250" is specific.
                // Rule 13: "if a player ends the game... scores -250 points".
                // This refers to the WINNER if they have Joker? "If a player ends the game while still holding a Joker...".
                // But winner discards last card. If they hold Joker, they can't end game unless Joker is set?
                // Ah, maybe they end game with Joker in hand (invalid win?).
                // But my logic ensures valid win.
                // Let's look at Rule 16: "All other players are scored... Joker: -250 points".
                // If 2-10 are 5 pts. 
                // It's ambiguous if 5 pts is penalty or reward.
                // Given "Rummy", usually leftover cards are penalty.
                // If Joker is -250, that's a massive reduction in penalty?
                // OR user means 250 points penalty.
                // I will stick to literal text:
                // Value: 2-10->5, JQK->10, A->15, Joker-> -250.
                
                score += -250
            } else {
                switch c.Rank {
                case models.Ace:
                    score += 15
                case models.Jack, models.Queen, models.King:
                    score += 10
                default:
                    score += 5
                }
            }

        }
        gm.Game.Players[i].Score += score // Accumulate score
    }
}

// RestartGame starts a new round without resetting scores
func (gm *GameManager) RestartGame(initiatorID string) error {
    gm.Mutex.Lock()
    defer gm.Mutex.Unlock()

    if gm.Game.Status != models.StateFinished {
        return errors.New("game is not finished")
    }

    if initiatorID != gm.Game.MasterPlayerID {
        return errors.New("only master can restart game")
    }

    // Reset Game State for new round
    gm.Game.Status = models.StateInProgress
    gm.Game.WinnerID = ""
    gm.Game.Pile = make([]models.Card, 0)
    gm.Game.TableSets = make([][]models.Card, 0) // Clear table sets

    // New Deck & Shuffle
    gm.Game.Deck = NewDeck()
    gm.Game.Deck = Shuffle(gm.Game.Deck)

    // Reset Players (Keep Score, ID, Name, Seat)
    for i := range gm.Game.Players {
        var handSize int
        if gm.Game.Players[i].ID == gm.Game.MasterPlayerID {
            handSize = 8
        } else {
            handSize = 7
        }

        hand := make([]models.Card, handSize)
        copy(hand, gm.Game.Deck[:handSize])
        gm.Game.Players[i].Hand = hand
        gm.Game.Deck = gm.Game.Deck[handSize:]
        
        // Reset round-specific flags
        gm.Game.Players[i].HasTakenFromPile = false
        gm.Game.Players[i].HasPlayedSet = false
        gm.Game.Players[i].PlayedSets = make([][]models.Card, 0)
    }

    // Master starts again
    masterIdx := 0
    for i, p := range gm.Game.Players {
        if p.ID == gm.Game.MasterPlayerID {
            masterIdx = i
            break
        }
    }
    gm.Game.CurrentTurnPlayer = masterIdx
    gm.Game.TurnPhase = models.PhasePlay // Master starts with 8 cards -> Play/Discard phase

    gm.save()
    return nil
}

func (gm *GameManager) GetPublicView(playerID string) *models.PublicGameView {
    view := &models.PublicGameView{
        ID: gm.Game.ID,
        Status: gm.Game.Status,
        CurrentTurnPlayer: gm.Game.CurrentTurnPlayer,
        TurnPhase: gm.Game.TurnPhase,
        DeckCount: len(gm.Game.Deck),
        Pile: gm.Game.Pile,
        TableSets: gm.Game.TableSets,
        OpponentHandSizes: make([]int, 4),
        WinnerID: gm.Game.WinnerID,
        MasterPlayerID: gm.Game.MasterPlayerID,
        Players: make([]models.PublicPlayer, 0),
    }

    for i, p := range gm.Game.Players {
        view.OpponentHandSizes[i] = len(p.Hand)
        
        view.Players = append(view.Players, models.PublicPlayer{
            ID: p.ID,
            Name: p.Name,
            SeatIndex: p.SeatIndex,
            Score: p.Score,
            IsConnected: p.IsConnected,
            HasPlayedSet: p.HasPlayedSet,
            PlayedSets: p.PlayedSets,
        })
        
        if p.ID == playerID {
            view.MyHand = p.Hand
            view.MySeatIndex = p.SeatIndex
            view.HasTakenFromPile = p.HasTakenFromPile
            view.HasPlayedSet = p.HasPlayedSet
            view.Score = p.Score
        }
    }
    return view
}

// Removing generic EndTurn since DiscardCard now handles turn switching.
