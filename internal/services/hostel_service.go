package services

import (
	"database/sql"
	"hostel-management/storage/repositories"
)

type HostelService interface {
	GetAllHostelNumbers() ([]int, error)
	GetHostelsInfo(db *sql.DB) ([]map[string]interface{}, error)
	GetHostelLocationByNumber(hostelNumber int) (string, error)
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

func (s *hostelService) GetHostelLocationByNumber(hostelNumber int) (string, error) {
	return s.hostelRepo.GetHostelLocationByNumber(hostelNumber)
}

func (s *hostelService) GetAllHostelNumbers() ([]int, error) {
	return s.hostelRepo.GetAllHostelNumbers()
}
