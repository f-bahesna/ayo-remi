// Package websocket provides the WebSocket infrastructure layer.
// It wraps the original ws package to expose clean interfaces.
package websocket

import (
	"remi-game/game"
	origws "remi-game/ws"
)

// Hub wraps the original WebSocket hub, managing client connections
// and message broadcasting.
type Hub struct {
	inner *origws.Hub
}

// NewHub creates and returns a new Hub wrapper.
func NewHub() *Hub {
	return &Hub{
		inner: origws.NewHub(),
	}
}

// Run starts the hub's main event loop. Should be called as a goroutine.
func (h *Hub) Run() {
	h.inner.Run()
}

// BroadcastGameUpdate sends personalized game state to all connected clients.
// Implements the usecase.Broadcaster interface.
func (h *Hub) BroadcastGameUpdate(gm *game.GameManager) {
	h.inner.BroadcastGameUpdate(gm)
}

// Inner returns the underlying ws.Hub for use with ServeWs.
func (h *Hub) Inner() *origws.Hub {
	return h.inner
}
