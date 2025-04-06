package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceService interface {
	CreateService(name, typeService, description string, is_date, is_hostel, is_phone bool, amount int) error
	GetAllServices() ([]models.Service, error)
	GetServiceByID(idInt int) (models.Service, error)
	UpdateServiceByID(idInt int, c *gin.Context) error
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
	service := models.Service{
		Name:        name,
		Type:        typeService,
		Amount:      amount,
		Description: description,
		Is_date:     is_date,
		Is_hostel:   is_hostel,
		Is_phone:    is_phone,
	}
	return s.repo.CreateService(service)
}

func (s *serviceService) GetAllServices() ([]models.Service, error) {
	services, err := s.repo.GetAllServices()
	if err != nil {
		return nil, err
	}
	for i := range services {
		services[i].Point = i + 1
	}
	return services, nil
}

func (s *serviceService) GetServiceByID(idInt int) (models.Service, error) {
	if idInt <= 0 {
		return models.Service{}, nil
	}
	return s.repo.GetServiceByID(idInt)
}

func (s *serviceService) UpdateServiceByID(idInt int, c *gin.Context) error {

	service := func(models.Service) models.Service {
		return models.Service{
			ID:          idInt,
			Name:        c.PostForm("name"),
			Type:        c.PostForm("type"),
			Description: c.PostForm("description"),
			Amount:      0,
			Is_date:     c.PostForm("is_date") == "on",
			Is_hostel:   c.PostForm("is_hostel") == "on",
			Is_phone:    c.PostForm("is_phone") == "on",
		}
	}(models.Service{})

	if amountStr := c.PostForm("amount"); amountStr != "" {
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return err
		}
		service.Amount = amount
	}

	return s.repo.UpdateServiceByID(idInt, service)
}

func (s *serviceService) DeleteService(idInt int) error {
	return s.repo.DeleteService(idInt)
}
