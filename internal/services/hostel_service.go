package services

import (
	"database/sql"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type HostelService interface {
	GetAllHostelNumbers() ([]int, error)
	GetHostelsInfo(db *sql.DB) ([]map[string]interface{}, error)
	GetHostelInfoByHeadman(db *sql.DB, email string) (map[string]interface{}, error)
	GetHostelInfo(hostelID int) (models.Hostel, error)
	InsertHeadmanIntoHostel(hostelID int, email string) error
	GetHostelLocationByNumber(hostelNumber int) (string, error)
	DeleteHeadmanFromHostel(hostelID int) error
}

type hostelService struct {
	hostelRepo repositories.HostelRepository
}

func NewHostelService(hostelRepo repositories.HostelRepository) HostelService {
	return &hostelService{
		hostelRepo: hostelRepo,
	}
}

func (s *hostelService) GetHostelsInfo(db *sql.DB) ([]map[string]interface{}, error) {
	hostels, err := s.hostelRepo.GetHostelsInfo(db)
	if err != nil {
		return nil, err
	}

	hostelData := []map[string]interface{}{}
	for _, hostel := range hostels {
		hostelData = append(hostelData, map[string]interface{}{
			"ID":             hostel.HostelID,
			"Number":         hostel.HostelNumber,
			"RoomsCount":     hostel.OccupiedRooms + hostel.AvailableRooms,
			"OccupiedRooms":  hostel.OccupiedRooms,
			"AvailableRooms": hostel.AvailableRooms,
			"ResidentsCount": hostel.ResidentsCount,
			"Location":       hostel.HostelLocation,
		})
	}

	return hostelData, nil
}

func (s *hostelService) GetHostelInfoByHeadman(db *sql.DB, email string) (map[string]interface{}, error) {
	hostel, err := s.hostelRepo.GetHostelInfoByHeadman(db, email)
	if err != nil {
		return nil, err
	}

	hostelData := map[string]interface{}{
		"ID":             hostel.HostelID,
		"Number":         hostel.HostelNumber,
		"RoomsCount":     hostel.OccupiedRooms + hostel.AvailableRooms,
		"OccupiedRooms":  hostel.OccupiedRooms,
		"AvailableRooms": hostel.AvailableRooms,
		"ResidentsCount": hostel.ResidentsCount,
		"Location":       hostel.HostelLocation,
	}

	return hostelData, nil
}

func (s *hostelService) GetHostelInfo(hostelID int) (models.Hostel, error) {
	hostel, err := s.hostelRepo.GetHostelInfo(hostelID)
	if err != nil {
		return models.Hostel{}, err
	}

	hostel.OccupiedPercent = int(float64(hostel.OccupiedRooms) / float64(hostel.RoomCount) * 100)
	hostel.AvailablePercent = int(float64(hostel.AvailableRooms) / float64(hostel.RoomCount) * 100)

	hostel.HostelContacts = hostel.HostelContacts[7:]
	hostel.HostelContacts = "8 " + hostel.HostelContacts

	return hostel, nil
}

func (s *hostelService) GetHostelLocationByNumber(hostelNumber int) (string, error) {
	return s.hostelRepo.GetHostelLocationByNumber(hostelNumber)
}

func (s *hostelService) GetAllHostelNumbers() ([]int, error) {
	return s.hostelRepo.GetAllHostelNumbers()
}

func (s *hostelService) InsertHeadmanIntoHostel(hostelID int, email string) error {
	return s.hostelRepo.AssignHeadmanToHostel(hostelID, email)
}

func (s *hostelService) DeleteHeadmanFromHostel(hostelID int) error {
	return s.hostelRepo.DeleteHeadmanFromHostel(hostelID, "")
}
