// Package dto defines Data Transfer Objects used at the application boundary.
// These types are used for communication between layers without exposing
// domain internals.
package dto

import "remi-game/models"

// PublicPlayer exposes safe player fields to all clients.
type PublicPlayer = models.PublicPlayer

// PublicGameView is the game state appropriate for a specific player,
// hiding other players' hands.
type PublicGameView = models.PublicGameView

// --- WebSocket Message Payloads ---

// DrawCardPayload specifies the source and count for a draw action.
type DrawCardPayload struct {
	Source string `json:"source"` // "DECK" or "PILE"
	Count  int    `json:"count"`
}

// DrawFromPilePayload specifies hand cards and pile card for a pile draw.
type DrawFromPilePayload struct {
	HandCardIDs []string `json:"handCardIds"`
	PileCardID  string   `json:"pileCardId"`
}

// PlaySetPayload contains the cards to play as a set.
type PlaySetPayload struct {
	Cards []models.Card `json:"cards"`
}

// DiscardCardPayload specifies which card to discard.
type DiscardCardPayload struct {
	CardID string `json:"cardId"`
}

// JoinGamePayload contains the player name for joining.
type JoinGamePayload struct {
	Name string `json:"name"`
}

// ErrorPayload contains an error message.
type ErrorPayload struct {
	Message string `json:"message"`
}
