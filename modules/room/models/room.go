package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SaveMessageReq struct {
		RoomID   primitive.ObjectID `json:"room_id"`
		Username string             `json:"username"`
		Content  string             `json:"content"`
	}
)
