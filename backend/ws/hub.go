package ws

import (
	"encoding/json"
	"remi-game/game"
	"sync"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

    mutex sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
            h.mutex.Lock()
			h.clients[client] = true
            h.mutex.Unlock()
		case client := <-h.unregister:
            h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
            h.mutex.Unlock()
		case message := <-h.broadcast:
            h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
            h.mutex.Unlock()
		}
	}
}

func (h *Hub) BroadcastGameUpdate(gm *game.GameManager) {
    h.mutex.Lock()
    defer h.mutex.Unlock()

    for client := range h.clients {
        // Generate personalized view
        view := gm.GetPublicView(client.PlayerID)
        
        msg := WSMessage{
            Type:    MsgGameUpdate, // Defined in messages.go
            Payload: view,
        }
        
        b, err := json.Marshal(msg)
        if err == nil {
            select {
            case client.Send <- b:
            default:
                close(client.Send)
                delete(h.clients, client)
            }
        }
    }
}
