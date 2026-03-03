// Package repository defines the port interfaces for persistence.
// These are implemented by infrastructure adapters (e.g., PostgreSQL).
package repository

import "remi-game/internal/domain/entity"

// GameRepository is the port for game persistence operations.
type GameRepository interface {
	// Save persists the current game state (upsert).
	Save(game entity.Game) error

	// FindByID retrieves a game by its unique identifier.
	FindByID(id string) (*entity.Game, error)
}
