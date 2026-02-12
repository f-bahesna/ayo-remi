package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"remi-game/game"
	"remi-game/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Client struct {
	Hub *Hub
	Conn *websocket.Conn
	Send chan []byte
    Game *game.GameManager
    PlayerID string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
        var msg WSMessage
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Printf("Invalid JSON: %v", err)
            continue
        }

        c.MapMessageToGameAction(msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, gameManager *game.GameManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256), Game: gameManager}
	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (c *Client) MapMessageToGameAction(msg WSMessage) {
    // Panic recovery: prevent crashing the WebSocket connection
    defer func() {
        if r := recover(); r != nil {
            log.Printf("PANIC in MapMessageToGameAction: %v (player: %s, action: %s)", r, c.PlayerID, msg.Type)
            c.SendError(fmt.Sprintf("internal error: %v", r))
        }
    }()

    // In a real app, use mutex on Game for all these calls
    c.Game.Mutex.Lock()
    defer c.Game.Mutex.Unlock()

    payloadBytes, _ := json.Marshal(msg.Payload)
    
    // For manual testing, let's log incoming messages
    log.Printf("Received: %s from %s", msg.Type, c.PlayerID)

    switch msg.Type {
    case MsgJoinGame:
        // We handle JoinGame specially because payload is struct map
        // Due to interface{}, json.Unmarshal might produce map[string]interface{}
        // We need robust unmarshalling.
        // Quick hack: Parse strictly if we know the structure.
        
        // Actually, msg.Payload is map[string]interface{} coming from ReadPump -> Unmarshal
        // So we marshal it back to bytes then strict unmarshal.
        var p JoinGamePayload
        json.Unmarshal(payloadBytes, &p)
        
        player, err := c.Game.AddPlayer(p.Name)
        if err != nil {
            c.SendError(err.Error())
            return
        }
        c.PlayerID = player.ID
        log.Printf("Player Joined: %s (%s)", p.Name, c.PlayerID)
        
        // Broadcast
        c.Hub.BroadcastGameUpdate(c.Game)

        if len(c.Game.Game.Players) == 4 && c.Game.Game.Status == models.StateWaiting {
            c.Game.StartGameUnlocked(c.Game.Game.MasterPlayerID)
            c.Hub.BroadcastGameUpdate(c.Game)
        }

    case MsgDrawCard:
        var p DrawCardPayload
        // Need to unmarshal payload to get source
        json.Unmarshal(payloadBytes, &p)
        source := p.Source
        if source == "" {
            source = "DECK"
        }
        
        if err := c.Game.DrawCard(c.PlayerID, source, p.Count); err != nil {
            c.SendError(err.Error())
            return
        }
        c.Hub.BroadcastGameUpdate(c.Game)

    case MsgDrawFromPile:
        var p DrawFromPilePayload
        json.Unmarshal(payloadBytes, &p)
        if err := c.Game.DrawFromPile(c.PlayerID, p.HandCardIDs, p.PileCardID); err != nil {
            c.SendError(err.Error())
            return
        }
        c.Hub.BroadcastGameUpdate(c.Game)

    case MsgPlaySet:
        var p PlaySetPayload
        json.Unmarshal(payloadBytes, &p)
        if err := c.Game.PlaySet(c.PlayerID, p.Cards); err != nil {
            c.SendError(err.Error())
            return
        }
        c.Hub.BroadcastGameUpdate(c.Game)

    case MsgDiscardCard:
        var p DiscardCardPayload
        json.Unmarshal(payloadBytes, &p)
        if err := c.Game.DiscardCard(c.PlayerID, p.CardID); err != nil {
            c.SendError(err.Error())
            return
        }
        c.Hub.BroadcastGameUpdate(c.Game)
        
    case MsgDeclareWin:
        if err := c.Game.DeclareWin(c.PlayerID); err != nil {
            c.SendError(err.Error())
            return
        }
        c.Hub.BroadcastGameUpdate(c.Game)

    case MsgStartGame:
        if err := c.Game.StartGame(c.PlayerID); err != nil {
             c.SendError(err.Error())
             return
        }
        c.Hub.BroadcastGameUpdate(c.Game)

    case MsgRestartGame:
        if err := c.Game.RestartGame(c.PlayerID); err != nil {
            c.SendError(err.Error())
            return
        }
        c.Hub.BroadcastGameUpdate(c.Game)
    }
}

func (c *Client) SendError(msg string) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("SendError recovered from panic: %v", r)
        }
    }()
    resp := WSMessage{Type: MsgError, Payload: ErrorPayload{Message: msg}}
    b, _ := json.Marshal(resp)
    select {
    case c.Send <- b:
    default:
        log.Printf("SendError: channel full or closed, dropping error: %s", msg)
    }
}
