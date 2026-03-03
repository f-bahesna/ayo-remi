package websocket

import (
	"errors"
	"sync"

	"remi-game/game"

	"github.com/google/uuid"
)

// Room represents a game room containing a game instance and its WebSocket hub.
type Room struct {
	ID   string
	Game *game.GameManager
	Hub  *Hub
}

// Registry manages all active game rooms.
type Registry struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

// NewRegistry creates a new room registry.
func NewRegistry() *Registry {
	return &Registry{
		rooms: make(map[string]*Room),
	}
}

// CreateRoom creates a new game room with its own GameManager and Hub.
// Returns the room ID.
func (r *Registry) CreateRoom() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()

	gm := game.NewGame(nil)
	gm.Game.ID = id

	hub := NewHub()
	go hub.Run()

	room := &Room{
		ID:   id,
		Game: gm,
		Hub:  hub,
	}

	r.rooms[id] = room
	return id
}

// GetRoom retrieves a room by ID.
func (r *Registry) GetRoom(id string) (*Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, ok := r.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}
