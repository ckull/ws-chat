package repository

import (
	"context"
	"fmt"
	"time"
	"ws-chat/modules/entities"
	"ws-chat/modules/room/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	roomRepository struct {
		db *mongo.Client
	}

	RoomRepository interface {
		messagesCollection() *mongo.Collection
		roomsCollection() *mongo.Collection
		// CreateRoom(room entities.Room) (primitive.ObjectID, error)
		SaveMessage(models.SaveMessageReq) error
	}
)

func NewRoomRepository(db *mongo.Client) RoomRepository {
	return &roomRepository{
		db,
	}
}

func (r *roomRepository) messagesCollection() *mongo.Collection {
	return r.db.Database("chat").Collection("messages")
}

func (r *roomRepository) roomsCollection() *mongo.Collection {
	return r.db.Database("chat").Collection("rooms")
}

func (r *roomRepository) CreateRoom(room entities.Room) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := r.roomsCollection()

	res, err := col.InsertOne(ctx, room)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create room: %v", err)
	}

	// Convert the InsertedID to a primitive.ObjectID type
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return insertedID, nil
}

func (r *roomRepository) SaveMessage(req models.SaveMessageReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := r.messagesCollection()

	// Add the `roomId` to the message before saving.
	msg := entities.Message{
		RoomID:    req.RoomID,
		Content:   req.Content,
		Username:  req.Username, // Assuming SaveMessageReq contains UserID and Content fields.
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := col.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}
	return nil

}
