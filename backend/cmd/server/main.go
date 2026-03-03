// Package main is the clean entry point for the Remi game server.
// It wires all layers together using dependency injection.
package main

import (
	"log"
	"net/http"

	infrahttp "remi-game/internal/infrastructure/http"
	infraws "remi-game/internal/infrastructure/websocket"
)

func main() {
	// --- Infrastructure Layer: WebSocket Registry ---
	registry := infraws.NewRegistry()

	// --- Infrastructure Layer: HTTP Handlers ---
	handler := infrahttp.NewHandler(registry)

	// --- Routes ---
	http.HandleFunc("/api/rooms", handler.CreateRoom)
	http.HandleFunc("/api/rooms/options", handler.CORSPreflight)
	http.HandleFunc("/ws", handler.WebSocket)

	// --- Start Server ---
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
