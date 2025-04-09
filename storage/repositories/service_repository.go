package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
)

type ServiceRepository interface {
	CreateService(service models.Service) error
	GetAllServices() ([]models.Service, error)
	GetServiceByID(idInt int) (models.Service, error)
	UpdateServiceByID(idInt int, service models.Service) error
	DeleteService(idInt int) error
}

type serviceRepository struct {
	db *sql.DB
}

func NewServiceRepository() ServiceRepository {
	return &serviceRepository{
		db: db.DB,
	}
}

func (r *serviceRepository) CreateService(service models.Service) error {
	query := "INSERT INTO Services (name, type, amount, description, is_date, is_hostel, is_phone) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, service.Name, service.Type, service.Amount, service.Description, service.Is_date, service.Is_hostel, service.Is_phone)
	if err != nil {
		return err
	}
	return nil
}

func (r *serviceRepository) GetAllServices() ([]models.Service, error) {
	query := "SELECT * FROM Services"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	services := []models.Service{}
	for rows.Next() {
		service := models.Service{}
		err := rows.Scan(&service.ID, &service.Name, &service.Type, &service.Amount, &service.Description, &service.Is_date, &service.Is_hostel, &service.Is_phone)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *serviceRepository) GetServiceByID(idInt int) (models.Service, error) {
	query := "SELECT * FROM Services WHERE id = ?"
	row := db.DB.QueryRow(query, idInt)

	service := models.Service{}
	err := row.Scan(&service.ID, &service.Name, &service.Type, &service.Amount, &service.Description, &service.Is_date, &service.Is_hostel, &service.Is_phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Service{}, fmt.Errorf("service with ID %d not found", idInt)
		}
		return models.Service{}, err
	}

	return service, nil
}

func (r *serviceRepository) UpdateServiceByID(idInt int, service models.Service) error {
	query := "UPDATE Services SET name = ?, type = ?, amount = ?, description = ?, is_date = ?, is_hostel = ?, is_phone = ? WHERE id = ?"
	_, err := db.DB.Exec(query, service.Name, service.Type, service.Amount, service.Description, service.Is_date, service.Is_hostel, service.Is_phone, idInt)
	if err != nil {
		return err
	}
	return nil
}

func (r *serviceRepository) DeleteService(idInt int) error {
	query := "DELETE FROM Services WHERE id = ?"
	_, err := db.DB.Exec(query, idInt)
	if err != nil {
		return err
	}
	return nil
}
