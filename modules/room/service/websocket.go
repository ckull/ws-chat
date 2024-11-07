package services

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
		RegisterClient(roomID string, client *entities.Client) error
		CreateRoom(name string, admin *entities.Client) *entities.Room
		Unregister(roomID string, client *entities.Client) error
	}
)

func NewWebSocketRoomService(roomRepository repository.RoomRepository) WebSocketRoomService {
	return &webSocketRoomService{
		Rooms:          make(map[string]*entities.Room),
		roomRepository: roomRepository,
	}
}

func (service *webSocketRoomService) CreateRoom(name string, admin *entities.Client) *entities.Room {
	roomID := primitive.NewObjectID().Hex()
	room := &entities.Room{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Admin:      admin,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Clients:    make(map[*websocket.Conn]*entities.Client),
		Broadcast:  make(chan *entities.Message),
		Register:   make(chan *entities.Client),
		Unregister: make(chan *entities.Client),
		Mu:         sync.Mutex{},
	}

	log.Printf("Room created, room id: %s", roomID)

	service.mu.Lock()
	service.Rooms[roomID] = room
	service.mu.Unlock()

	go service.Run(roomID)

	return room
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

func (service *webSocketRoomService) RegisterClient(roomID string, client *entities.Client) error {
	service.mu.Lock()
	room, exists := service.Rooms[roomID]
	service.mu.Unlock()

	if !exists {
		return errors.New("Room ID doesn't exist")
	}

	room.Register <- client

	return nil
}

func (service *webSocketRoomService) Unregister(roomID string, client *entities.Client) error {
	service.mu.Lock()
	room, exists := service.Rooms[roomID]
	service.mu.Unlock()

	if !exists {
		return errors.New("Room ID doesn't exist")
	}

	room.Unregister <- client

	return nil
}
