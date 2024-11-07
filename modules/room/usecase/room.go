package usecase

import (
	"ws-chat/modules/room/models"
	"ws-chat/modules/room/repository"
)

type (
	roomUsecase struct {
		repository repository.RoomRepository
	}

	RoomUsecase interface {
		SaveMessage(req models.SaveMessageReq) error
	}
)

func NewRoomUsecase(repository repository.RoomRepository) RoomUsecase {
	return &roomUsecase{
		repository,
	}
}

func (u *roomUsecase) SaveMessage(req models.SaveMessageReq) error {

	return u.repository.SaveMessage(req)
}
