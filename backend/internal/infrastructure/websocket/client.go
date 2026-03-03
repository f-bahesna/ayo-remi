package websocket

import (
	"net/http"
	"remi-game/game"
	origws "remi-game/ws"
)

// ServeWs upgrades an HTTP connection to WebSocket and registers the client.
// This delegates to the original ws.ServeWs function.
func ServeWs(hub *Hub, gameManager *game.GameManager, w http.ResponseWriter, r *http.Request) {
	origws.ServeWs(hub.Inner(), gameManager, w, r)
}
