package models

type Suit string
const (
    Spades   Suit = "S"
    Hearts   Suit = "H"
    Diamonds Suit = "D"
    Clubs    Suit = "C"
    Joker    Suit = "J"
)

type Rank int
// 1-9 represent numeric values directly
// 11=J, 12=Q, 13=K, 14=A
// 0=Joker
const (
	JokerRank Rank = 0
	Jack      Rank = 11
	Queen     Rank = 12
	King      Rank = 13
	Ace       Rank = 14
)

type Card struct {
    Suit Suit `json:"suit"`
    Rank Rank `json:"rank"`
    ID   string `json:"id"` // Unique ID to track instances
}

type TurnPhase string
const (
	PhaseDraw    TurnPhase = "DRAW"
	PhasePlay    TurnPhase = "PLAY" // Includes playing sets
	PhaseDiscard TurnPhase = "DISCARD"
)

type Player struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Hand             []Card   `json:"hand"`
	SeatIndex        int      `json:"seatIndex"` // 0-3
	IsConnected      bool     `json:"isConnected"`
	Score             int       `json:"score"`
    HasTakenFromPile  bool      `json:"hasTakenFromPile"` // "Center Black Card Rule": once per game
    HasPlayedSet      bool      `json:"hasPlayedSet"`     // Track if player has melded at least one set
    PlayedSets        [][]Card  `json:"playedSets"`       // Sets played by this player
}

type GameState string

const (
	StateWaiting    GameState = "WAITING"
	StateInProgress GameState = "IN_PROGRESS"
	StateFinished   GameState = "FINISHED"
)

type Game struct {
	ID                string    `json:"id"`
	Status            GameState `json:"status"`
	CurrentTurnPlayer int       `json:"currentTurnPlayer"` // SeatIndex of current player
	TurnPhase         TurnPhase `json:"turnPhase"`
	Deck              []Card    `json:"deck"`              // Remaining cards
	Pile              []Card    `json:"pile"`              // Center pile (Discard pile)
	TableSets         [][]Card  `json:"tableSets"`         // Public sets on table
	Players           []Player  `json:"players"`
	WinnerID          string    `json:"winnerId,omitempty"`
    MasterPlayerID    string    `json:"masterPlayerId"`    // ID of the room creator/master
}

// PublicPlayer exposes safe fields to everyone
type PublicPlayer struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    SeatIndex    int      `json:"seatIndex"`
    Score        int      `json:"score"`
    IsConnected  bool     `json:"isConnected"`
    HasPlayedSet bool     `json:"hasPlayedSet"`
    PlayedSets   [][]Card `json:"playedSets"`
}

// PublicView returns a game state appropriate for a specific player (hiding other hands)
type PublicGameView struct {
	ID                string         `json:"id"`
	Status            GameState      `json:"status"`
	CurrentTurnPlayer int            `json:"currentTurnPlayer"`
	TurnPhase         TurnPhase      `json:"turnPhase"`
	DeckCount         int            `json:"deckCount"`
	Pile              []Card         `json:"pile"`
	TableSets         [][]Card       `json:"tableSets"`
	MyHand            []Card         `json:"myHand"`
	MySeatIndex       int            `json:"mySeatIndex"`
	OpponentHandSizes []int          `json:"opponentHandSizes"` // Ordered by seat index 0-3
	WinnerID          string         `json:"winnerId,omitempty"`
	HasTakenFromPile  bool           `json:"hasTakenFromPile"`
    MasterPlayerID    string         `json:"masterPlayerId"`
    
    // Extra fields for UI
    Players           []PublicPlayer `json:"players"`
    HasPlayedSet      bool           `json:"hasPlayedSet"`     // Current player's status
    Score             int            `json:"score"`            // Current player's score
}
