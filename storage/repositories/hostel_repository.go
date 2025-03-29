package repositories

import (
	"database/sql"
	"hostel-management/storage/db"
	"hostel-management/storage/models"
)

type HostelRepository interface {
	GetAllHostelNumbers() ([]int, error)
	GetHostelsInfo(db *sql.DB) ([]models.Hostel, error)
	GetHostelLocationByNumber(hostelNumber int) (string, error)
}

type hostelRepository struct {
	db *sql.DB
}

func NewHostelRepository() HostelRepository {
	return &hostelRepository{
		db: db.DB,
	}
}

func (r *hostelRepository) GetAllHostelNumbers() ([]int, error) {
	rows, err := db.DB.Query("SELECT number FROM Hostels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hostelNumbers []int
	for rows.Next() {
		var hostelNumber int
		if err := rows.Scan(&hostelNumber); err != nil {
			return nil, err
		}
		hostelNumbers = append(hostelNumbers, hostelNumber)
	}

	return hostelNumbers, nil
}

func (r *hostelRepository) GetHostelsInfo(db *sql.DB) ([]models.Hostel, error) {
	query := `
        SELECT 
            h.id AS hostel_id,
            h.number AS hostel_number,
			h.location AS hostel_location,
            COUNT(DISTINCT u.id) AS residents_count,
            SUM(CASE WHEN r.status = 'occupied' THEN 1 ELSE 0 END) AS occupied_rooms,
            SUM(CASE WHEN r.status = 'unoccupied' THEN 1 ELSE 0 END) AS available_rooms
        FROM 
            Hostels h
        LEFT JOIN 
            Rooms r ON r.Hostels_id = h.id
        LEFT JOIN 
            Users u ON u.Rooms_id = r.id
        GROUP BY 
            h.id, h.number;
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hostels []models.Hostel
	for rows.Next() {
		var hostel models.Hostel
		if err := rows.Scan(
			&hostel.HostelID,
			&hostel.HostelNumber,
			&hostel.HostelLocation,
			&hostel.ResidentsCount,
			&hostel.OccupiedRooms,
			&hostel.AvailableRooms,
		); err != nil {
			return nil, err
		}
		hostels = append(hostels, hostel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hostels, nil
}

func (r *hostelRepository) GetHostelLocationByNumber(hostelNumber int) (string, error) {
	var location string
	err := db.DB.QueryRow("SELECT location FROM Hostels WHERE number = ?", hostelNumber).Scan(&location)
	if err != nil {
		return "", err
	}
	return location, nil
}
