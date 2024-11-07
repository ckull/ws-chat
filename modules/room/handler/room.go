package handler

import (
	"net/http"
	"ws-chat/modules/room/models"
	"ws-chat/modules/room/usecase"

	"github.com/labstack/echo/v4"
)

type (
	roomHandler struct {
		roomUsecase usecase.RoomUsecase
	}

	RoomHandler interface {
		SaveMessage(c echo.Context) error
	}
)

func NewRoomHandler(roomUsecase usecase.RoomUsecase) RoomHandler {
	return &roomHandler{
		roomUsecase,
	}
}

func (h *roomHandler) SaveMessage(c echo.Context) error {
	var saveMessageReq models.SaveMessageReq

	if err := c.Bind(&saveMessageReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := h.roomUsecase.SaveMessage(saveMessageReq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Success"})
}
