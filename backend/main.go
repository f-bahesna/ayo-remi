package main

import (
    "log"
    "net/http"
    "remi-game/ws"
)

func main() {
    // Initialize Registry (manages Rooms: Game + Hub)
    registry := ws.NewRegistry()

    // API: Create Room
    http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        roomID := registry.CreateRoom()
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(`{"roomId": "` + roomID + `"}`))
    })
    
    // CORS Preflight
    http.HandleFunc("/api/rooms/options", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.WriteHeader(http.StatusOK)
    })

    // WebSocket Endpoint
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        // Query params: ?room=ID&name=Name
        roomID := r.URL.Query().Get("room")
        if roomID == "" {
            http.Error(w, "Missing room ID", http.StatusBadRequest)
            return
        }

        room, err := registry.GetRoom(roomID)
        if err != nil {
            http.Error(w, "Room not found", http.StatusNotFound)
            return
        }

        ws.ServeWs(room.Hub, room.Game, w, r)
    })

    log.Println("Server started on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
