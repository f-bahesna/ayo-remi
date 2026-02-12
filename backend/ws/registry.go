package ws

import (
    "errors"
    "sync"
    "remi-game/game"
    "github.com/google/uuid"
)

type Room struct {
    ID   string
    Game *game.GameManager
    Hub  *Hub
}

type Registry struct {
    mu    sync.RWMutex
    rooms map[string]*Room
}

func NewRegistry() *Registry {
    return &Registry{
        rooms: make(map[string]*Room),
    }
}

func (r *Registry) CreateRoom() string {
    r.mu.Lock()
    defer r.mu.Unlock()

    id := uuid.New().String()
    
    // Create Game and Hub for this room
    gm := game.NewGame(nil)
    gm.Game.ID = id
    
    hub := NewHub()
    go hub.Run() // Start the hub
    
    room := &Room{
        ID:   id,
        Game: gm,
        Hub:  hub,
    }

    r.rooms[id] = room
    return id
}

func (r *Registry) GetRoom(id string) (*Room, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    room, ok := r.rooms[id]
    if !ok {
        return nil, errors.New("room not found")
    }
    return room, nil
}
