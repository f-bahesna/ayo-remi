// Package service contains domain services that encapsulate
// business rules not naturally belonging to a single entity.
package service

import (
	"remi-game/game"
	"remi-game/internal/domain/entity"
)

// CardService encapsulates card-related domain logic:
// deck creation, shuffling, set validation, and win checking.
type CardService struct{}

// NewCardService creates a new CardService.
func NewCardService() *CardService {
	return &CardService{}
}

// NewDeck creates a standard 54-card deck (52 + 2 Jokers).
func (s *CardService) NewDeck() []entity.Card {
	return game.NewDeck()
}

// Shuffle returns a shuffled copy of the given deck.
func (s *CardService) Shuffle(deck []entity.Card) []entity.Card {
	return game.Shuffle(deck)
}

// IsValidSet checks whether a group of cards forms a valid run or set.
// Rules: Run = same suit, consecutive ranks (min 3). Set = same rank, different suits (min 3).
// Jokers act as wildcards.
func (s *CardService) IsValidSet(cards []entity.Card) bool {
	return game.IsValidSet(cards)
}

// CanFormSetWith checks if adding newCard to hand can form at least one valid set.
func (s *CardService) CanFormSetWith(hand []entity.Card, newCard entity.Card) bool {
	return game.CanFormSetWith(hand, newCard)
}

// CanFormSetWithMultiple checks if ALL newCards can be part of a single valid set
// when combined with cards from the hand.
func (s *CardService) CanFormSetWithMultiple(hand []entity.Card, newCards []entity.Card) bool {
	return game.CanFormSetWithMultiple(hand, newCards)
}

// IsWinningHand checks if the hand can be fully formed into valid sets
// with at most 1 card remaining (the discard).
func (s *CardService) IsWinningHand(hand []entity.Card) bool {
	return game.IsWinningHand(hand)
}
