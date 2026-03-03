// Package http provides HTTP handlers for the REST API endpoints.
package http

import (
	"net/http"

	infraws "remi-game/internal/infrastructure/websocket"
)

// Handler holds dependencies for HTTP request handling.
type Handler struct {
	Registry *infraws.Registry
}

// NewHandler creates a new Handler with the given registry.
func NewHandler(registry *infraws.Registry) *Handler {
	return &Handler{Registry: registry}
}

// CreateRoom handles POST /api/rooms — creates a new game room.
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := h.Registry.CreateRoom()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(`{"roomId": "` + roomID + `"}`))
}

// CORSPreflight handles OPTIONS /api/rooms/options — CORS preflight.
func (h *Handler) CORSPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.WriteHeader(http.StatusOK)
}

// WebSocket handles GET /ws — upgrades to WebSocket connection.
func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	room, err := h.Registry.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	infraws.ServeWs(room.Hub, room.Game, w, r)
}
