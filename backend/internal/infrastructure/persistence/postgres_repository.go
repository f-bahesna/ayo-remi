// Package persistence provides infrastructure adapters for game persistence.
// Implements the domain repository interface using PostgreSQL via the original db package.
package persistence

import (
	"remi-game/db"
	"remi-game/internal/domain/entity"
)

// PostgresGameRepository implements repository.GameRepository
// by delegating to the original db.GameStore.
type PostgresGameRepository struct {
	store *db.GameStore
}

// NewPostgresGameRepository creates a new PostgresGameRepository.
func NewPostgresGameRepository(store *db.GameStore) *PostgresGameRepository {
	return &PostgresGameRepository{store: store}
}

// Save persists the game state to PostgreSQL.
func (r *PostgresGameRepository) Save(game entity.Game) error {
	return r.store.SaveGame(game)
}

// FindByID retrieves a game by ID from PostgreSQL.
func (r *PostgresGameRepository) FindByID(id string) (*entity.Game, error) {
	return r.store.LoadGame(id)
}
