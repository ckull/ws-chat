package handler

import (
	"fmt"
	"net/http"
	"time"
	"ws-chat/modules/room/service/ws"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	roomService ws.WebSocketRoomService
}

var (
	upgrader = websocket.Upgrader{}
)

func NewWebSocketHandler(roomService ws.WebSocketRoomService) *WebSocketHandler {
	return &WebSocketHandler{roomService: roomService}
}

func (h *WebSocketHandler) CreateRoom(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upgrade to WebSocket"})
	}

	roomName := c.QueryParam("room")
	host := c.QueryParam("host")

	client := ws.NewClient(host, conn)

	room := h.roomService.CreateRoom(roomName, client)

	go client.ReadPump()

	msg := fmt.Sprintf(`{"success": "Create room succeeded, %s"}`, room.ID)

	return conn.WriteMessage(websocket.TextMessage, []byte(msg))

}

func (h *WebSocketHandler) JoinRoom(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upgrade to WebSocket"})
	}

	userName := c.QueryParam("username")

	roomID := c.QueryParam("room_id")

	if userName == "" || roomID == "" {
		conn.Close() // Ensure connection is closed if input validation fails
		errMsg := `{"error": "Username and Room ID are required"}`
		return conn.WriteMessage(websocket.TextMessage, []byte(errMsg))
	}

	client := &ws.Client{Conn: conn, Username: userName, Send: make(chan []byte), JoinedAt: time.Now()}

	err = h.roomService.RegisterClient(roomID, client)
	if err != nil {

		return conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}

	msg := `{"success": "Join room succeed"}`

	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
