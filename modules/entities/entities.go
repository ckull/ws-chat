package entities

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// User represents a user entity that is stored in MongoDB and used in JSON responses.
	User struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB ObjectID
		Username  string             `bson:"username" json:"username"`          // User's nickname
		FirstName string             `bson:"firstName" json:"firstName"`        // User's first name
		LastName  string             `bson:"lastName" json:"lastName"`          // User's last name
		Email     string             `bson:"email" json:"email"`                // User's email address
		UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`       // Timestamp for the last update
		CreatedAt time.Time          `bson:"created_at" json:"createdAt"`       // Timestamp for creation
	}

	Client struct {
		ID       string
		Username string
		Conn     *websocket.Conn
		JoinedAt time.Time
		Send     chan []byte
	}

	// Room represents a chat room entity.
	Room struct {
		ID         primitive.ObjectID          `bson:"_id,omitempty" json:"id,omitempty"` // Room identifier
		Name       string                      `bson:"name" json:"name"`                  // Room name
		Admin      *Client                     `bson:"admin" json:"admin"`                // Admin's WebSocket connection (not stored)
		CreatedAt  time.Time                   `bson:"created_at" json:"createdAt"`       // Timestamp for creation
		UpdatedAt  time.Time                   `bson:"updated_at" json:"updatedAt"`       // Timestamp for the last update
		Clients    map[*websocket.Conn]*Client // Map of WebSocket connections to User (not stored in MongoDB or JSON)
		Broadcast  chan *Message               // Channel for broadcasting messages (not stored)
		Register   chan *Client
		Unregister chan *Client
		Mu         sync.Mutex // Mutex for thread-safe operations (not stored)
	}

	// Message represents a message sent in a WebSocket session.
	Message struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // Room identifier
		RoomID    primitive.ObjectID `bson:"room_id" json:"room_id"`
		Username  string             `bson:"username" json:"username"`    // Reference to the User who sent the message
		Content   string             `bson:"content" json:"content"`      // Message content
		UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"` // Timestamp for the last update
		CreatedAt time.Time          `bson:"created_at" json:"createdAt"` // Timestamp for creation
	}
)
