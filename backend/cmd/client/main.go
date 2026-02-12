package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func connectAndPlay(name string, done chan struct{}) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Handle incoming messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("%s received: %s", name, message)
		}
	}()

	// Join Game
	joinMsg := WSMessage{
		Type: "JOIN_GAME",
		Payload: map[string]string{"name": name},
	}
	b, _ := json.Marshal(joinMsg)
	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("write:", err)
		return
	}

	// Wait a bit
	time.Sleep(2 * time.Second)

    // Keep connection open
    select {}
}

func main() {
    numPlayers := flag.Int("n", 4, "number of players to simulate")
    flag.Parse()
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    // Simulate N players
    for i := 1; i <= *numPlayers; i++ {
        go connectAndPlay(string("Bot"+string(rune('0'+i))), make(chan struct{}))
        time.Sleep(500 * time.Millisecond)
    }

    <-interrupt
}
