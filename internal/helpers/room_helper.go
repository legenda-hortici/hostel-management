package helpers

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type RoomHelper struct {
	roomService repositories.RoomRepository
	userService repositories.UserRepository
}

func NewRoomHelper() RoomHelper {
	return RoomHelper{
		roomService: repositories.NewRoomRepository(),
		userService: repositories.NewUserRepository(),
	}
}

// Функция для перевода типа комнаты на русский
func TranslateRoomType(roomType string) string {
	switch roomType {
	case "once":
		return "Одноместная"
	case "double":
		return "Двухместная"
	case "triple":
		return "Трехместная"
	case "premium double":
		return "Двухместная (комфорт)"
	case "premium triple":
		return "Трехместная (комфорт)"
	default:
		return "Неизвестный тип"
	}
}

// Функция для перевода статуса комнаты на русский
func TranslateRoomStatus(status string) string {
	switch status {
	case "unoccupied":
		return "Доступна"
	case "occupied":
		return "Занята"
	case "renovation":
		return "На ремонте"
	default:
		return "Неизвестный статус"
	}
}

// Функция для перевода информации о жильцах на русский
func TranslateUserCount(userCount int) int {
	// Оставляем только само число
	return userCount
}

func (rh RoomHelper) TranslateRoom(room models.Room) models.Room {
	room.Type = TranslateRoomType(room.Type)
	room.Status = TranslateRoomStatus(room.Status)
	room.UserCount = TranslateUserCount(room.UserCount)

	return room
}

func (rh RoomHelper) ValidateRoomData(room models.Room) error {
	if room.Number == 0 || room.Type == "" || room.Status == "" || room.HostelNumber == 0 {
		return errors.New("ValidateRoomData: заполните все обязательные поля")
	}

	_, err := rh.roomService.CreateRoom(room.Type, room.Status, room.Number, room.UserCount, room.HostelNumber)
	if err != nil {
		return errors.New("ValidateResidentData: ошибка при создании комнаты")
	}

	return nil
}

func (rh RoomHelper) ValidateAddResidentData(email string, room int) error {
	if email == "" {
		return errors.New("ValidateResidentEmail: заполните обязательное поле")
	}

	user, err := rh.userService.GetByEmail(email)
	if err != nil {
		return errors.New("ValidateResidentEmail: пользователь не найден")
	}

	err = rh.roomService.InsertResidentIntoRoom(room, user.Email)
	if err != nil {
		return errors.New("ValidateResidentEmail: ошибка при добавлении жильца в комнату")
	}

	return nil
}

func (rh RoomHelper) ValidateDeleteResidentData(email string) (int, error) {
	if email == "" {
		return 0, errors.New("ValidateResidentEmail: заполните обязательное поле")
	}

	user, err := rh.userService.GetByEmail(email)
	if err != nil {
		return 0, errors.New("ValidateResidentEmail: пользователь не найден")
	}

	// Получаем текущую комнату пользователя
	roomId, err := rh.roomService.GetRoomIdByEmail(user.Email)
	if err != nil {
		return 0, errors.New("ValidateResidentEmail: ошибка при получении комнаты пользователя")
	}

	// Вызываем удаление жильца
	err = rh.roomService.DeleteResidentFromRoom(user.Email)
	if err != nil {
		return 0, errors.New("ValidateResidentEmail: ошибка при удалении жильца из комнаты")
	}

	// Возвращаем ID комнаты, из которой был удален жилец
	return roomId, nil
}

func (rh RoomHelper) DefineRoomStatus(roomType string, userCount int, roomNumber int) (string, error) {
	roomCapacity := map[string]int{
		"once":           1,
		"double":         2,
		"triple":         3,
		"premium double": 2,
		"premium triple": 3,
	}

	capacity, exists := roomCapacity[roomType]
	if !exists {
		return "", errors.New("DefineRoomStatus: неизвестный тип комнаты")
	}

	if userCount == capacity {
		err := rh.roomService.UpdateRoomStatus(roomNumber, "occupied")
		if err != nil {
			return "", errors.New("DefineRoomStatus: ошибка при обновлении статуса комнаты")
		}
		return "occupied", nil
	}
	err := rh.roomService.UpdateRoomStatus(roomNumber, "unoccupied")
	if err != nil {
		return "", errors.New("DefineRoomStatus: ошибка при обновлении статуса комнаты")
	}
	return "unoccupied", nil
}
