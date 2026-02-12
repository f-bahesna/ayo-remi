package game

import (
    "remi-game/models"
    "testing"
)

// MockStore for testing
type MockStore struct{}

func (m *MockStore) SaveGame(game models.Game) error {
    return nil
}

func TestStartGame_DealCounts(t *testing.T) {
    gm := NewGame(&MockStore{})
    
    // Add 4 players
    gm.AddPlayer("P1")
    gm.AddPlayer("P2")
    gm.AddPlayer("P3")
    gm.AddPlayer("P4")
    
    err := gm.StartGame(gm.Game.MasterPlayerID)
    if err != nil {
        t.Fatalf("StartGame failed: %v", err)
    }
    
    // Check P1 (Master) has 8 cards
    if len(gm.Game.Players[0].Hand) != 8 {
        t.Errorf("Master should have 8 cards, got %d", len(gm.Game.Players[0].Hand))
    }
    
    // Check others have 7
    for i := 1; i < 4; i++ {
        if len(gm.Game.Players[i].Hand) != 7 {
            t.Errorf("Player %d should have 7 cards, got %d", i, len(gm.Game.Players[i].Hand))
        }
    }
    
    // Check Deck size
    // 52 + 2 jokers = 54 cards total.
    // Dealt: 8 + 7*3 = 29.
    // Pile: 0 (initially).
    // Deck should be 54 - 29 = 25.
    if len(gm.Game.Deck) != 25 {
        t.Errorf("Deck should have 25 cards, got %d", len(gm.Game.Deck))
    }
}

func TestDrawCard_PileAndReshuffle(t *testing.T) {
    gm := NewGame(&MockStore{})
    gm.AddPlayer("Master")
    gm.AddPlayer("P2")
    gm.StartGame(gm.Game.MasterPlayerID)
    
    // Master starts. 
    // Phase is PLAY (because 8 cards).
    // Let's force phase to DRAW for testing, or simulate a turn.
    // Master Discard to end turn.
    
    master := gm.Game.Players[0]
    discardCard := master.Hand[0]
    // Simulate discard logic manually or use DiscardCard if available? 
    // State.go doesn't have DiscardCard in the snippets I saw, but it was mentioned.
    // I'll manually manipulate state for unit test speed.
    
    gm.Game.Pile = append(gm.Game.Pile, discardCard)
    gm.Game.Players[0].Hand = master.Hand[1:]
    gm.Game.TurnPhase = models.PhaseDraw
    gm.Game.CurrentTurnPlayer = 1 // P2's turn
    
    // P2 Draws from Pile
    // FIX: Manually set P2 hand and Pile to ensure valid set (Restriction compliance)
    validPileCard := models.Card{Suit: models.Spades, Rank: 5, ID: "pile1"}
    gm.Game.Pile = []models.Card{validPileCard}
    
    // P2 has 4S, 6S (Run with 5S)
    gm.Game.Players[1].Hand = []models.Card{
        {Suit: models.Spades, Rank: 4, ID: "h1"},
        {Suit: models.Spades, Rank: 6, ID: "h2"},
        {Suit: models.Hearts, Rank: 10, ID: "h3"}, // filler
        {Suit: models.Hearts, Rank: 10, ID: "h4"},
        {Suit: models.Hearts, Rank: 10, ID: "h5"},
        {Suit: models.Hearts, Rank: 10, ID: "h6"},
        {Suit: models.Hearts, Rank: 10, ID: "h7"},
    }
    
    initialPileSize := len(gm.Game.Pile)
    err := gm.DrawCard(gm.Game.Players[1].ID, "PILE", 1)
    if err != nil {
        t.Fatalf("Draw from Pile failed: %v", err)
    }
    
    if len(gm.Game.Pile) != initialPileSize - 1 {
        t.Errorf("Pile should decrease by 1")
    }
    if len(gm.Game.Players[1].Hand) != 8 { // 7+1
        t.Errorf("P2 should have 8 cards after draw")
    }
    
    // Test Reshuffle (Deck is empty, Pile has cards)
    gm.Game.Deck = []models.Card{}
    
    // Pile needs > 1 card to reshuffle. Top card stays. Rest goes to deck.
    // Let's add 2 cards to pile.
    topCard := models.Card{Suit: models.Hearts, Rank: 1, ID: "top"}
    bottomCard := models.Card{Suit: models.Hearts, Rank: 2, ID: "bottom"}
    gm.Game.Pile = []models.Card{bottomCard, topCard}
    
    // Switch to Master turn
    gm.Game.TurnPhase = models.PhaseDraw
    gm.Game.CurrentTurnPlayer = 0 // Master
    
    // Master draws from DECK (empty -> triggers reshuffle)
    err = gm.DrawCard(master.ID, "DECK", 1)
    if err != nil {
        t.Fatalf("Draw from empty Deck (reshuffle) failed: %v", err)
    }
    
    // Verify results
    // - Pile should have 1 card (topCard)
    // - Deck should have 0 (bottomCard moved to deck, then drawn by Master)
    // - Master hand +1
    
    if len(gm.Game.Pile) != 1 {
          t.Errorf("Pile should have 1 card left")
    }
    if gm.Game.Pile[0].ID != topCard.ID {
        t.Errorf("Top card of pile should persist")
    }
}

func TestDrawPileRestriction(t *testing.T) {
	// Setup generic game
    store := &MockStore{}
    gm := NewGame(store)
    p1 := models.Player{ID: "p1", Name: "P1", Hand: []models.Card{
		{Suit: models.Spades, Rank: 2, ID: "h1"},
		{Suit: models.Hearts, Rank: 10, ID: "h2"}, // Random cards, no set potential with 3S
	}}
	gm.Game = models.Game{
        ID: "test",
        Status: models.StateInProgress,
        Players: []models.Player{p1},
        CurrentTurnPlayer: 0,
        TurnPhase: models.PhaseDraw,
        Pile: []models.Card{{Suit: models.Spades, Rank: 3, ID: "p1"}}, // 3S
    }

	// 1. Try to draw (Should FAIL: 2s, 10h + 3s is not a set)
	err := gm.DrawCard("p1", "PILE", 1)
	if err == nil {
		t.Errorf("Expected error when drawing from pile without set potential")
	}

	// 2. Add Spades 4 to hand (now 2s, 4s. + 3s = Run of 2-3-4 Spades)
	gm.Game.Players[0].Hand = append(gm.Game.Players[0].Hand, models.Card{Suit: models.Spades, Rank: 4, ID: "h3"})
	
	// Try again (Should SUCCEED)
	err = gm.DrawCard("p1", "PILE", 1)
	if err != nil {
		t.Errorf("Expected success drawing from pile with set potential: %v", err)
	}
    
    if len(gm.Game.Players[0].Hand) != 4 { // 2s, 10h, 4s + 3s
        t.Errorf("Hand size check failed")
    }
}

func TestDrawMultipleFromPile(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)
    
    // Pile: [6S, 7S, 8S] (Bottom to Top)
    // Hand: [5S]
    // Draw 3 (6,7,8) + 5S = 5,6,7,8 Run -> Valid
    
    p1 := models.Player{ID: "p1", Name: "P1", Hand: []models.Card{
        {Suit: models.Spades, Rank: 5, ID: "h1"},
    }}
    
    gm.Game = models.Game{
        ID: "test",
        Status: models.StateInProgress,
        Players: []models.Player{p1},
        CurrentTurnPlayer: 0,
        TurnPhase: models.PhaseDraw,
        Pile: []models.Card{
            {Suit: models.Spades, Rank: 6, ID: "p1"},
            {Suit: models.Spades, Rank: 7, ID: "p2"},
            {Suit: models.Spades, Rank: 8, ID: "p3"},
        }, 
    }
    
    // Draw 3
    err := gm.DrawCard("p1", "PILE", 3)
    if err != nil {
        t.Fatalf("Failed to draw 3 valid cards: %v", err)
    }
    
    if len(gm.Game.Players[0].Hand) != 4 {
        t.Errorf("Should have 4 cards")
    }
    if len(gm.Game.Pile) != 0 {
        t.Errorf("Pile should be empty")
    }
}

func TestDrawFromPile(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)

    // Hand: [4S, 6S, 10H, JH, QH]
    // Pile: [5S, 9D, KH]
    // Pick 5S from pile with 4S+6S from hand → Run 4-5-6 Spades ✓

    p1 := models.Player{ID: "p1", Name: "P1", Hand: []models.Card{
        {Suit: models.Spades, Rank: 4, ID: "h1"},
        {Suit: models.Spades, Rank: 6, ID: "h2"},
        {Suit: models.Hearts, Rank: 10, ID: "h3"},
        {Suit: models.Hearts, Rank: 11, ID: "h4"},
        {Suit: models.Hearts, Rank: 12, ID: "h5"},
    }}

    gm.Game = models.Game{
        ID:     "test",
        Status: models.StateInProgress,
        Players: []models.Player{p1},
        CurrentTurnPlayer: 0,
        TurnPhase: models.PhaseDraw,
        Pile: []models.Card{
            {Suit: models.Spades, Rank: 5, ID: "pile1"},
            {Suit: models.Diamonds, Rank: 9, ID: "pile2"},
            {Suit: models.Hearts, Rank: 13, ID: "pile3"},
        },
        TableSets: make([][]models.Card, 0),
    }

    // 1. Invalid: cards don't form a set (4S + 10H + 5S = no)
    err := gm.DrawFromPile("p1", []string{"h1", "h3"}, "pile1")
    if err == nil {
        t.Errorf("Expected error for invalid set")
    }
    // Reset phase since error shouldn't change state
    gm.Game.TurnPhase = models.PhaseDraw

    // 2. Invalid: pile card not found
    err = gm.DrawFromPile("p1", []string{"h1", "h2"}, "nonexistent")
    if err == nil {
        t.Errorf("Expected error for missing pile card")
    }
    gm.Game.TurnPhase = models.PhaseDraw

    // 3. Invalid: not enough hand cards
    err = gm.DrawFromPile("p1", []string{"h1"}, "pile1")
    if err == nil {
        t.Errorf("Expected error for too few hand cards")
    }
    gm.Game.TurnPhase = models.PhaseDraw

    // 4. Valid: 4S + 6S + 5S = Run
    err = gm.DrawFromPile("p1", []string{"h1", "h2"}, "pile1")
    if err != nil {
        t.Fatalf("DrawFromPile failed: %v", err)
    }

    // Hand should have 3 cards (5 - 2 used)
    if len(gm.Game.Players[0].Hand) != 3 {
        t.Errorf("Hand should have 3 cards, got %d", len(gm.Game.Players[0].Hand))
    }

    // Pile should have 2 cards (3 - 1 picked)
    if len(gm.Game.Pile) != 2 {
        t.Errorf("Pile should have 2 cards, got %d", len(gm.Game.Pile))
    }

    // Set should be played to table
    if len(gm.Game.Players[0].PlayedSets) != 1 {
        t.Errorf("Should have 1 played set, got %d", len(gm.Game.Players[0].PlayedSets))
    }

    // Phase should be PLAY
    if gm.Game.TurnPhase != models.PhasePlay {
        t.Errorf("Phase should be PLAY, got %s", gm.Game.TurnPhase)
    }
}
