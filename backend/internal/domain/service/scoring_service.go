package service

import "remi-game/internal/domain/entity"

// ScoringService encapsulates the scoring rules for end-of-game calculation.
// Card point values: 2-10 → 5, J/Q/K → 10, A → 15, Joker → -250.
type ScoringService struct{}

// NewScoringService creates a new ScoringService.
func NewScoringService() *ScoringService {
	return &ScoringService{}
}

// CalculateCardScore returns the point value of a single card.
func (s *ScoringService) CalculateCardScore(card entity.Card) int {
	if card.Suit == entity.Joker {
		return -250
	}
	switch card.Rank {
	case entity.Ace:
		return 15
	case entity.Jack, entity.Queen, entity.King:
		return 10
	default:
		return 5
	}
}

// CalculateHandScore returns the total score for a player's remaining hand.
func (s *ScoringService) CalculateHandScore(hand []entity.Card) int {
	score := 0
	for _, c := range hand {
		score += s.CalculateCardScore(c)
	}
	return score
}
