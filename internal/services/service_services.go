package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type ServiceService interface {
	CreateService(name, typeService, description string, is_date, is_hostel, is_phone bool, amount int) error
	GetAllServices() ([]models.Service, error)
	GetServiceByID(idInt int) (models.Service, error)
	UpdateServiceByID(idInt int, service models.Service) error
	DeleteService(idInt int) error
}

type serviceService struct {
	repo repositories.ServiceRepository
}

func NewServiceService(repo repositories.ServiceRepository) ServiceService {
	return &serviceService{
		repo: repo,
	}
}

func (s *serviceService) CreateService(name, typeService, description string, is_date, is_hostel, is_phone bool, amount int) error {
	return s.repo.CreateService(name, typeService, description, is_date, is_hostel, is_phone, amount)
}

func (s *serviceService) GetAllServices() ([]models.Service, error) {
	return s.repo.GetAllServices()
}

func (s *serviceService) GetServiceByID(idInt int) (models.Service, error) {
	return s.repo.GetServiceByID(idInt)
}

func (s *serviceService) UpdateServiceByID(idInt int, service models.Service) error {
	return s.repo.UpdateServiceByID(idInt, service)
}

func (s *serviceService) DeleteService(idInt int) error {
	return s.repo.DeleteService(idInt)
}
