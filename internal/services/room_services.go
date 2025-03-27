package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type RoomService interface {
	CreateRoom(roomType, status string, number, userCount, hostelNumber int) (int64, error)
	GetAllRooms() ([]models.Room, error)
	GetHostelIDByName(hostelNumber int) (int, error)
	GetRoomByID(roomID int) (models.Room, error)
	GetRoomNumberByID(roomID int) (int, error)
	GetRoomByNumber(roomNumber int) (*models.Room, error)
	GetRoomID(roomNumber string) (int, error)
	GetResidentsByRoomID(roomID int) ([]models.User, error)
	InsertResidentIntoRoom(roomID int, email string) error
	DeleteResidentFromRoom(email string) error
	GetRoomIdByEmail(email string) (int, error)
	GetRoomsCount() (int, error)
	GetHostelIDByRoomID(roomID int) (int, error)
	GetHostelNumberByID(hostelID int) (int, error)
	UpdateRoomStatus(roomID int, status string) error
	FreezeRoom(roomID int) error
	GetInventoryByRoomID(id int) ([]models.Inventory, error)
	GetRoomNumberByRoomID(roomID int) (int, error)
}

type roomServiceImpl struct {
	roomRepo repositories.RoomRepository
}

func NewRoomService(roomRepo repositories.RoomRepository) RoomService {
	return &roomServiceImpl{roomRepo: roomRepo}
}

func (s *roomServiceImpl) CreateRoom(roomType, status string, number, userCount, hostelNumber int) (int64, error) {
	return s.roomRepo.CreateRoom(roomType, status, number, userCount, hostelNumber)
}

func (s *roomServiceImpl) GetAllRooms() ([]models.Room, error) {

	return s.roomRepo.GetAllRooms()
}

func (s *roomServiceImpl) GetHostelIDByName(hostelNumber int) (int, error) {
	return s.roomRepo.GetHostelIDByName(hostelNumber)
}

func (s *roomServiceImpl) GetRoomByID(roomID int) (models.Room, error) {
	return s.roomRepo.GetRoomByID(roomID)
}

func (s *roomServiceImpl) GetRoomNumberByID(roomID int) (int, error) {
	return s.roomRepo.GetRoomNumberByID(roomID)
}

func (s *roomServiceImpl) GetRoomByNumber(roomNumber int) (*models.Room, error) {
	return s.roomRepo.GetRoomByNumber(roomNumber)
}

func (s *roomServiceImpl) GetRoomID(roomNumber string) (int, error) {
	return s.roomRepo.GetRoomID(roomNumber)
}

func (s *roomServiceImpl) GetResidentsByRoomID(roomID int) ([]models.User, error) {
	return s.roomRepo.GetResidentsByRoomID(roomID)
}

func (s *roomServiceImpl) InsertResidentIntoRoom(roomID int, email string) error {
	return s.roomRepo.InsertResidentIntoRoom(roomID, email)
}

func (s *roomServiceImpl) DeleteResidentFromRoom(email string) error {
	return s.roomRepo.DeleteResidentFromRoom(email)
}

func (s *roomServiceImpl) GetRoomIdByEmail(email string) (int, error) {
	return s.roomRepo.GetRoomIdByEmail(email)
}

func (s *roomServiceImpl) GetRoomsCount() (int, error) {
	return s.roomRepo.GetRoomsCount()
}

func (s *roomServiceImpl) GetHostelIDByRoomID(roomID int) (int, error) {
	return s.roomRepo.GetHostelIDByRoomID(roomID)
}

func (s *roomServiceImpl) GetHostelNumberByID(hostelID int) (int, error) {
	return s.roomRepo.GetHostelNumberByID(hostelID)
}

func (s *roomServiceImpl) UpdateRoomStatus(roomID int, status string) error {
	return s.roomRepo.UpdateRoomStatus(roomID, status)
}

func (s *roomServiceImpl) FreezeRoom(roomID int) error {
	return s.roomRepo.FreezeRoom(roomID)
}

func (s *roomServiceImpl) GetInventoryByRoomID(id int) ([]models.Inventory, error) {
	return s.roomRepo.GetInventoryByRoomID(id)
}

func (s *roomServiceImpl) GetRoomNumberByRoomID(roomID int) (int, error) {
	return s.roomRepo.GetRoomNumberByRoomID(roomID)
}
