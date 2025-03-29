package services

import (
	"database/sql"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type HostelService interface {
	GetAllHostelNumbers() ([]int, error)
	GetHostelsInfo(db *sql.DB) ([]models.Hostel, error)
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

func (s *hostelService) GetHostelsInfo(db *sql.DB) ([]models.Hostel, error) {
	return s.hostelRepo.GetHostelsInfo(db)
}

func (s *hostelService) GetHostelLocationByNumber(hostelNumber int) (string, error) {
	return s.hostelRepo.GetHostelLocationByNumber(hostelNumber)
}

func (s *hostelService) GetAllHostelNumbers() ([]int, error) {
	return s.hostelRepo.GetAllHostelNumbers()
}
