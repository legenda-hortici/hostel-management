package services

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type RoomService interface {
	CreateRoom(roomType, status string, number, userCount, hostelNumber int) error
	GetAllRooms() ([]models.Room, error)
	GetHostelIDByName(hostelNumber int) (int, error)
	GetRoomByID(roomID int) (models.Room, error)
	GetRoomNumberByID(roomID int) (int, error)
	GetRoomByNumber(roomNumber int) (*models.Room, error)
	GetRoomID(roomNumber string) (int, error)
	GetResidentsByRoomID(roomID int) ([]models.User, error)
	InsertResidentIntoRoom(roomID int, email string) error
	DeleteResidentFromRoom(email string) (int, error)
	GetRoomIdByEmail(email string) (int, error)
	GetRoomsCount() (int, error)
	GetHostelIDByRoomID(roomID int) (int, error)
	GetHostelNumberByID(hostelID int) (int, error)
	UpdateRoomStatus(roomID int, status string) error
	FreezeRoom(roomID int) error
	UnfreezeRoom(roomID int) error
	GetInventoryByRoomID(id int) ([]models.Inventory, error)
	GetRoomNumberByRoomID(roomID int) (int, error)
	
}

type roomServiceImpl struct {
	roomRepo repositories.RoomRepository
}

func NewRoomService(roomRepo repositories.RoomRepository) RoomService {
	return &roomServiceImpl{roomRepo: roomRepo}
}

func (s *roomServiceImpl) DefineRoomStatus(roomType string, userCount int, roomNumber int) (string, error) {
	roomCapacity := map[string]int{
		"одноместная":           1,
		"двухместная":           2,
		"трехместная":           3,
		"двухместная (премиум)": 2,
		"трехместная (премиум)": 3,
	}

	capacity, exists := roomCapacity[roomType]
	if !exists {
		return "", errors.New("DefineRoomStatus: неизвестный тип комнаты")
	}

	if userCount == capacity {
		err := s.UpdateRoomStatus(roomNumber, "Занята")
		if err != nil {
			return "", errors.New("DefineRoomStatus: ошибка при обновлении статуса комнаты")
		}
		return "Занята", nil
	}
	err := s.UpdateRoomStatus(roomNumber, "Доступна")
	if err != nil {
		return "", errors.New("DefineRoomStatus: ошибка при обновлении статуса комнаты")
	}
	return "Доступна", nil
}

func (s *roomServiceImpl) CreateRoom(roomType, status string, number, userCount, hostelNumber int) error {
	if number == 0 || roomType == "" || status == "" || hostelNumber == 0 {
		return errors.New("fields cannot be empty")
	}

	validTypes := map[string]bool{
		"одноместная":           true,
		"двухместная":           true,
		"трехместная":           true,
		"двухместная (премиум)": true,
		"трехместная (премиум)": true,
	}

	validStatuses := map[string]bool{
		"Доступна":   true,
		"Занята":     true,
		"На ремонте": true,
	}

	if !validTypes[roomType] {
		return errors.New("invalid room type")
	}

	if !validStatuses[status] {
		return errors.New("invalid room status")
	}

	room := models.Room{
		Number:       number,
		Type:         roomType,
		Status:       status,
		HostelNumber: hostelNumber,
		UserCount:    0,
	}

	return s.roomRepo.CreateRoom(room)
}

func (s *roomServiceImpl) GetAllRooms() ([]models.Room, error) {
	rooms, err := s.roomRepo.GetAllRooms()
	if err != nil {
		return nil, err
	}
	for room := range rooms {
		if rooms[room].Status != "На ремонте" {
			s.DefineRoomStatus(rooms[room].Type, rooms[room].UserCount, rooms[room].Number)
		}
	}
	return rooms, nil
}

// Функция для получения ID хостела по его номеру
func (s *roomServiceImpl) GetHostelIDByName(hostelNumber int) (int, error) {
	return s.roomRepo.GetHostelIDByName(hostelNumber)
}

// Функция для получения комнаты по ее ID
func (s *roomServiceImpl) GetRoomByID(roomID int) (models.Room, error) {
	room, err := s.roomRepo.GetRoomByID(roomID)
	if err != nil {
		return models.Room{}, err
	}

	residents, err := s.GetResidentsByRoomID(roomID)
	if err != nil {
		return models.Room{}, err
	}
	room.UserCount = len(residents)

	return room, nil
}

// Функция для получения номера комнаты по ее ID
func (s *roomServiceImpl) GetRoomNumberByID(roomID int) (int, error) {
	return s.roomRepo.GetRoomNumberByID(roomID)
}

// Функция для получения комнаты по ее номеру
func (s *roomServiceImpl) GetRoomByNumber(roomNumber int) (*models.Room, error) {
	return s.roomRepo.GetRoomByNumber(roomNumber)
}

// Функция для получения ID комнаты по ее номеру
func (s *roomServiceImpl) GetRoomID(roomNumber string) (int, error) {
	return s.roomRepo.GetRoomID(roomNumber)
}

// Функция для получения жильцов в комнате
func (s *roomServiceImpl) GetResidentsByRoomID(roomID int) ([]models.User, error) {
	return s.roomRepo.GetResidentsByRoomID(roomID)
}

// Функция для добавления жильца в комнату
func (s *roomServiceImpl) InsertResidentIntoRoom(roomID int, email string) error {
	if email == "" {
		return errors.New("ValidateResidentEmail: заполните обязательное поле")
	}
	return s.roomRepo.InsertResidentIntoRoom(roomID, email)
}

// Функция для удаления жильца из комнаты
func (s *roomServiceImpl) DeleteResidentFromRoom(email string) (int, error) {
	if email == "" {
		return 0, errors.New("ValidateResidentEmail: заполните обязательное поле")
	}
	roomId, err := s.GetRoomIdByEmail(email)
	if err != nil {
		return 0, errors.New("ValidateResidentEmail: ошибка при получении комнаты пользователя")
	}
	return roomId, s.roomRepo.DeleteResidentFromRoom(email)
}

// Функция для получения ID комнаты по ее email
func (s *roomServiceImpl) GetRoomIdByEmail(email string) (int, error) {
	return s.roomRepo.GetRoomIdByEmail(email)
}

// Функция для получения количества комнат
func (s *roomServiceImpl) GetRoomsCount() (int, error) {
	return s.roomRepo.GetRoomsCount()
}

// Функция для получения ID хостела по ID комнаты
func (s *roomServiceImpl) GetHostelIDByRoomID(roomID int) (int, error) {
	return s.roomRepo.GetHostelIDByRoomID(roomID)
}

// Функция для получения номера хостела по его ID
func (s *roomServiceImpl) GetHostelNumberByID(hostelID int) (int, error) {
	return s.roomRepo.GetHostelNumberByID(hostelID)
}

// Функция для обновления статуса комнаты
func (s *roomServiceImpl) UpdateRoomStatus(roomID int, status string) error {
	return s.roomRepo.UpdateRoomStatus(roomID, status)
}

// Функция для заморозки комнаты
func (s *roomServiceImpl) FreezeRoom(roomID int) error {
	return s.roomRepo.FreezeRoom(roomID)
}

// Функция для получения инвентаря комнаты
func (s *roomServiceImpl) GetInventoryByRoomID(id int) ([]models.Inventory, error) {
	return s.roomRepo.GetInventoryByRoomID(id)
}

// Функция для получения номера комнаты по ее ID
func (s *roomServiceImpl) GetRoomNumberByRoomID(roomID int) (int, error) {
	return s.roomRepo.GetRoomNumberByRoomID(roomID)
}

func (s *roomServiceImpl) UnfreezeRoom(roomID int) error {
	return s.roomRepo.UnfreezeRoom(roomID)
}