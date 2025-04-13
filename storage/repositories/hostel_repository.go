package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
)

type HostelRepository interface {
	GetAllHostelNumbers() ([]int, error)
	GetHostelsInfo(db *sql.DB) ([]models.Hostel, error)
	GetHostelInfo(id int) (models.Hostel, error)
	GetHostelLocationByNumber(hostelNumber int) (string, error)
	AssignHeadmanToHostel(hostelID int, email string) error
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

			-- Общее количество жильцов (через подзапрос)
			(
				SELECT COUNT(*) 
				FROM Users u
				JOIN Rooms r2 ON u.Rooms_id = r2.id
				WHERE r2.Hostels_id = h.id
			) AS residents_count,

			-- Количество занятых комнат
			SUM(CASE WHEN r.status = 'Занята' AND r.id != 999 THEN 1 ELSE 0 END) AS occupied_rooms,

			-- Количество доступных комнат (с учетом статуса "Доступна" или "На ремонте")
			SUM(CASE WHEN r.status IN ('Доступна', 'На ремонте') AND r.id != 999 THEN 1 ELSE 0 END) AS available_rooms

		FROM 
			Hostels h

		-- JOIN только с Rooms (без дублирования)
		LEFT JOIN Rooms r ON r.Hostels_id = h.id

		GROUP BY 
			h.id, h.number, h.location;
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

func (r *hostelRepository) AssignHeadmanToHostel(hostelID int, email string) error {
	var userExists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM Users WHERE email = ?)"
	err := db.DB.QueryRow(queryCheck, email).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке email: %v", err)
	}

	if !userExists {
		return fmt.Errorf("комендант с таким email не найден: %s", email)
	}

	queryUpdate := "UPDATE Hostels SET headman_email = ? WHERE id = ?"
	_, err = db.DB.Exec(queryUpdate, email, hostelID)
	if err != nil {
		return fmt.Errorf("ошибка при назначении коменданта: %v", err)
	}

	return nil
}

func (r *hostelRepository) GetHostelInfo(id int) (models.Hostel, error) {

	query := `
			SELECT 
				h.id AS hostel_id,
				h.number AS hostel_number,
				h.location AS hostel_location,
				h.contacts AS hostel_contacts,
				h.room_count AS hostel_room_count,

				-- Староста: ищем по email, указанному в общежитии
				hu.name AS headman_name,
				hu.surname AS headman_surname,
				hu.email AS headman_email,

				-- Кол-во жильцов, связанных с комнатами этого общежития
				(
					SELECT COUNT(*) 
					FROM Users u
					JOIN Rooms r ON u.Rooms_id = r.id
					WHERE r.Hostels_id = h.id AND r.id != 999
				) AS residents_count,

				-- Комнаты по статусам
				SUM(CASE WHEN r.status = 'Занята' THEN 1 ELSE 0 END) AS occupied_rooms,
				SUM(CASE WHEN r.status = 'Доступна' THEN 1 ELSE 0 END) AS available_rooms

			FROM 
				Hostels h

			-- JOIN всех комнат этого общежития (для подсчета статусов)
			LEFT JOIN Rooms r ON r.Hostels_id = h.id AND r.id != 999

			-- JOIN старосты по email
			LEFT JOIN Users hu ON hu.email = h.headman_email

			WHERE 
				h.id = ?

			GROUP BY 
				h.id, h.number, h.location, h.contacts, h.room_count,
				hu.name, hu.surname, hu.email

			`

	row := r.db.QueryRow(query, id)

	hostel := &models.Hostel{}
	err := row.Scan(&hostel.HostelID,
		&hostel.HostelNumber,
		&hostel.HostelLocation,
		&hostel.HostelContacts,
		&hostel.RoomCount,
		&hostel.HeadmanName,
		&hostel.HeadmanSurname,
		&hostel.HeadmanEmail,
		&hostel.ResidentsCount,
		&hostel.OccupiedRooms,
		&hostel.AvailableRooms)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Hostel{}, nil
		}
		return models.Hostel{}, err
	}

	return *hostel, nil
}
