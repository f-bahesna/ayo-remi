// Package usecase contains application use cases that orchestrate
// domain logic and infrastructure. Each use case method represents
// a single user action.
package usecase

import (
	"remi-game/game"
	"remi-game/internal/application/dto"
	"remi-game/internal/domain/entity"
	"remi-game/internal/domain/repository"
	"remi-game/models"
)

// Broadcaster defines the interface for pushing game updates to clients.
// Implemented by the WebSocket infrastructure layer.
type Broadcaster interface {
	BroadcastGameUpdate(gm *game.GameManager)
}

// GameUseCase orchestrates game operations by coordinating
// between the domain layer (GameManager) and infrastructure (repository, broadcaster).
type GameUseCase struct {
	repo        repository.GameRepository
	broadcaster Broadcaster
}

// NewGameUseCase creates a new GameUseCase with the given dependencies.
func NewGameUseCase(repo repository.GameRepository, broadcaster Broadcaster) *GameUseCase {
	return &GameUseCase{
		repo:        repo,
		broadcaster: broadcaster,
	}
}

// AddPlayer adds a player to the game and broadcasts the update.
func (uc *GameUseCase) AddPlayer(gm *game.GameManager, name string) (*entity.Player, error) {
	player, err := gm.AddPlayer(name)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// StartGame initiates the game (only master can start).
func (uc *GameUseCase) StartGame(gm *game.GameManager, initiatorID string) error {
	return gm.StartGame(initiatorID)
}

// StartGameUnlocked initiates the game without acquiring lock (caller must hold lock).
func (uc *GameUseCase) StartGameUnlocked(gm *game.GameManager, initiatorID string) error {
	return gm.StartGameUnlocked(initiatorID)
}

// DrawCard draws a card from the specified source.
func (uc *GameUseCase) DrawCard(gm *game.GameManager, playerID string, source string, count int) error {
	return gm.DrawCard(playerID, source, count)
}

// DrawFromPile picks a specific card from the pile using hand cards to form a set.
func (uc *GameUseCase) DrawFromPile(gm *game.GameManager, playerID string, handCardIDs []string, pileCardID string) error {
	return gm.DrawFromPile(playerID, handCardIDs, pileCardID)
}

// PlaySet plays a valid set of cards from the player's hand to the table.
func (uc *GameUseCase) PlaySet(gm *game.GameManager, playerID string, cards []entity.Card) error {
	return gm.PlaySet(playerID, cards)
}

// DiscardCard discards a card and ends the player's turn.
func (uc *GameUseCase) DiscardCard(gm *game.GameManager, playerID string, cardID string) error {
	return gm.DiscardCard(playerID, cardID)
}

// DeclareWin declares the current player as the winner.
func (uc *GameUseCase) DeclareWin(gm *game.GameManager, playerID string) error {
	return gm.DeclareWin(playerID)
}

// RestartGame starts a new round without resetting scores.
func (uc *GameUseCase) RestartGame(gm *game.GameManager, initiatorID string) error {
	return gm.RestartGame(initiatorID)
}

// GetPublicView returns a player-specific view of the game state.
func (uc *GameUseCase) GetPublicView(gm *game.GameManager, playerID string) *dto.PublicGameView {
	return gm.GetPublicView(playerID)
}

// IsAutoStartReady checks if the game should auto-start (4 players joined, still waiting).
func (uc *GameUseCase) IsAutoStartReady(gm *game.GameManager) bool {
	return len(gm.Game.Players) == 4 && gm.Game.Status == models.StateWaiting
}

// GetMasterPlayerID returns the master player ID for auto-start.
func (uc *GameUseCase) GetMasterPlayerID(gm *game.GameManager) string {
	return gm.Game.MasterPlayerID
}
