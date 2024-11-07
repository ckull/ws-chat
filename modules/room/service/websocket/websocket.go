package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"ws-chat/modules/entities"
	"ws-chat/modules/room/repository"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	webSocketRoomService struct {
		Rooms          map[string]*entities.Room
		roomRepository repository.RoomRepository
		mu             sync.Mutex // Protects access to the rooms map
	}

	WebSocketRoomService interface {
		Run(roomID string)
		RegisterClient(roomID string, client *Client) error
		CreateRoom(name string, admin *Client) *entities.Room
		Unregister(roomID string, client *Client) error
	}
)

func NewWebSocketRoomService(roomRepository repository.RoomRepository) WebSocketRoomService {
	return &webSocketRoomService{
		Rooms:          make(map[string]*entities.Room),
		roomRepository: roomRepository,
	}
}

func (s *webSocketRoomService) CreateRoom(name string, admin *Client) *entities.Room {
	roomID := primitive.NewObjectID().Hex()
	room := &entities.Room{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Admin:      admin,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Clients:    make(map[*websocket.Conn]*Client),
		Broadcast:  make(chan *entities.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Mu:         sync.Mutex{},
	}

	log.Printf("Room created, room id: %s", roomID)

	s.mu.Lock()
	s.Rooms[roomID] = room
	s.mu.Unlock()

	go s.Run(roomID)

	return room
}

func (s *webSocketRoomService) HandleWebSocket(conn *websocket.Conn, roomID string, username string) error {
	// Create a new client instance
	client := NewClient(username, username, conn)

	// Register the client with the room
	err := s.RegisterClient(roomID, client)
	if err != nil {
		conn.Close()
		return err
	}

	// Start the ReadPump and WritePump as goroutines
	go client.ReadPump(s.Rooms[roomID]) // Start reading messages from client and forwarding to Broadcast
	// go client.WritePump() // Start sending messages from Send channel to client

	return nil
}

func (service *webSocketRoomService) Run(roomID string) {
	service.mu.Lock()
	room, exists := service.Rooms[roomID]
	service.mu.Unlock()

	if !exists {
		return
	}

	log.Printf("Room created, room id: %s", roomID)

	for {
		select {
		case client := <-room.Register:
			room.Mu.Lock()
			room.Clients[client.Conn] = client
			room.Mu.Unlock()
			log.Printf("Client %s joined room %s\n", client.Username, room.Name)

		case client := <-room.Unregister:
			room.Mu.Lock()
			if _, ok := room.Clients[client.Conn]; ok {
				delete(room.Clients, client.Conn)
				close(client.Send)
				log.Printf("Client %s left room %s\n", client.Username, room.Name)
			}
			room.Mu.Unlock()

		case message := <-room.Broadcast:
			room.Mu.Lock()

			broadcastMsg, err := json.Marshal(message)
			if err != nil {
				log.Printf("Failed to marshal message: %v\n", err)
				continue
			}

			for conn, client := range room.Clients {
				select {
				case client.Send <- broadcastMsg:
					if err := conn.WriteMessage(websocket.TextMessage, broadcastMsg); err != nil {
						fmt.Printf("Error sending message to %s: %v\n", client.Username, err)
						conn.Close()
						delete(room.Clients, conn)
					}
				default:
					close(client.Send)
					delete(room.Clients, conn)
					conn.Close()
					log.Printf("Client %s disconnected from room %s\n", client.Username, room.Name)
				}
			}
			room.Mu.Unlock()
		}
	}
}

func (s *webSocketRoomService) RegisterClient(roomID string, client *Client) error {
	s.mu.Lock()
	room, exists := s.Rooms[roomID]
	s.mu.Unlock()

	if !exists {
		return errors.New("Room ID doesn't exist")
	}

	room.Register <- client

	return nil
}

func (s *webSocketRoomService) Unregister(roomID string, client *Client) error {
	s.mu.Lock()
	room, exists := s.Rooms[roomID]
	s.mu.Unlock()

	if !exists {
		return errors.New("Room ID doesn't exist")
	}

	room.Unregister <- client

	return nil
}
