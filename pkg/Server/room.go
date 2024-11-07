package server

import (
	"ws-chat/modules/room/handler"
	"ws-chat/modules/room/repository"
	services "ws-chat/modules/room/service"
)

func (s *Server) roomService() {

	roomRepo := repository.NewRoomRepository(s.Db)
	wsRoomService := services.NewWebSocketRoomService(roomRepo)

	wsRoomHandler := handler.NewWebSocketHandler(wsRoomService)

	s.App.GET("/room/create", wsRoomHandler.CreateRoom)
	s.App.GET("/room/join", wsRoomHandler.JoinRoom)

}
