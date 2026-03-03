// Package entity defines the core domain types for the Remi card game.
// These are re-exports from the original models package, establishing
// the canonical import path for domain entities in Clean Architecture.
package entity

import "remi-game/models"

// --- Value Objects ---

// Suit represents a card suit (Spades, Hearts, Diamonds, Clubs, Joker).
type Suit = models.Suit

const (
	Spades   = models.Spades
	Hearts   = models.Hearts
	Diamonds = models.Diamonds
	Clubs    = models.Clubs
	Joker    = models.Joker
)

// Rank represents a card rank (2-10, J, Q, K, A, Joker).
type Rank = models.Rank

const (
	JokerRank = models.JokerRank
	Jack      = models.Jack
	Queen     = models.Queen
	King      = models.King
	Ace       = models.Ace
)

// TurnPhase represents the current phase of a player's turn.
type TurnPhase = models.TurnPhase

const (
	PhaseDraw    = models.PhaseDraw
	PhasePlay    = models.PhasePlay
	PhaseDiscard = models.PhaseDiscard
)

// GameState represents the overall state of a game.
type GameState = models.GameState

const (
	StateWaiting    = models.StateWaiting
	StateInProgress = models.StateInProgress
	StateFinished   = models.StateFinished
)

// --- Entities ---

// Card represents a playing card with a unique ID.
type Card = models.Card

// Player represents a player in the game.
type Player = models.Player

// Game represents the full game state (Aggregate Root).
type Game = models.Game
