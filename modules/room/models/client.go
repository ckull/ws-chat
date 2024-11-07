package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type (
	Client struct {
		Username string
		Conn     *websocket.Conn
		JoinedAt time.Time
		Send     chan []byte
	}
)
