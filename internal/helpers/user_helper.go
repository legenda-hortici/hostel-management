package helpers

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
	"strconv"
)

type UserHelper interface {
	ValidateResidentData(resident models.User) error
}

type userHelper struct {
	roomService repositories.RoomRepository
	userService repositories.UserRepository
}

func NewUserHelper() UserHelper {
	return &userHelper{
		roomService: repositories.NewRoomRepository(),
		userService: repositories.NewUserRepository(),
	}
}

func (uh *userHelper) ValidateResidentData(resident models.User) error {
	if resident.Username == "" || resident.Email == "" || resident.Password == "" || resident.Institute.String == "" || resident.RoomNumber == 0 {
		return errors.New("ValidateResidentData: заполните все обязательные поля")
	}
	roomNumber := strconv.Itoa(resident.RoomNumber)

	roomID, err := uh.roomService.GetRoomID(roomNumber)
	if err != nil {
		return errors.New("ValidateResidentData: комната не найдена")
	}
	resident.Room_id = roomID
	resident.Role = "user"

	err = uh.userService.Create(&resident)
	if err != nil {
		return errors.New("ValidateResidentData: ошибка при создании пользователя")
	}

	return nil
}
