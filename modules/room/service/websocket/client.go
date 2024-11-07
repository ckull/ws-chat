package websocket

import (
	"encoding/json"
	"log"
	"time"
	"ws-chat/modules/entities"

	"github.com/gorilla/websocket"
)

type (
	Client struct {
		ID       string
		Username string
		Conn     *websocket.Conn
		JoinedAt time.Time
		Send     chan []byte
	}
)

func NewClient(ID string, Username string, conn *websocket.Conn) *Client {
	return &Client{
		ID:       ID,
		Username: Username,
		Conn:     conn,
		JoinedAt: time.Now(),
		Send:     make(chan []byte), // Creates a buffered channel to hold outgoing messages
	}
}

func (client *Client) ReadPump(room *entities.Room) {

	for {
		// Read a message from the WebSocket connection
		_, messageData, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", client.Username, err)
			break
		}

		// Unmarshal the incoming message if needed (e.g., if it's JSON)
		var msg entities.Message
		err = json.Unmarshal(messageData, &msg)
		if err != nil {
			log.Printf("Error unmarshaling message from %s: %v", client.Username, err)
			continue
		}

		// Forward the message to the room's Broadcast channel
		room.Broadcast <- &msg
	}
}
