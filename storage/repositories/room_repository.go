package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
	"log"
)

type RoomRepository interface {
	CreateRoom(room models.Room) error
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
	GetRoomNumberByRoomID(userID int) (int, error)
}

type roomRepository struct {
	db            *sql.DB
	inventoryRepo InventoryRepository
}

func NewRoomRepository() RoomRepository {
	return &roomRepository{
		db:            db.DB,
		inventoryRepo: NewInventoryRepository(),
	}
}

func (r *roomRepository) CreateRoom(room models.Room) error {
	hostelID, err := r.GetHostelIDByName(room.HostelNumber)
	if err != nil {
		return fmt.Errorf("ошибка при получении ID общежития: %w", err)
	}

	query := "INSERT INTO Rooms (type, status, number, user_count, Hostels_id) VALUES (?, ?, ?, ?, ?)"
	_, err = r.db.Exec(query, room.Type, room.Status, room.Number, room.UserCount, hostelID)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении комнаты: %w", err)
	}

	return nil
}

func (r *roomRepository) GetAllRooms() ([]models.Room, error) {
	query := `SELECT r.id, r.type, r.status, r.number, r.user_count, h.number
				FROM Rooms r
				JOIN Hostels h ON r.Hostels_id = h.id;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := []models.Room{}
	for rows.Next() {
		room := models.Room{}
		err := rows.Scan(&room.ID, &room.Type, &room.Status, &room.Number, &room.UserCount, &room.HostelNumber)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *roomRepository) GetHostelIDByName(hostelNumber int) (int, error) {
	var hostelID int
	query := "SELECT id FROM Hostels WHERE number = ? LIMIT 1"
	err := r.db.QueryRow(query, hostelNumber).Scan(&hostelID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении ID общежития: %w", err)
	}
	return hostelID, nil
}

func (r *roomRepository) GetRoomByID(roomID int) (models.Room, error) {
	query := `
		SELECT r.id, r.type, r.status, r.number, 
			LEAST(COUNT(u.id), 
				CASE 
					WHEN r.type = 'одноместная' THEN 1 
					WHEN r.type = 'двухместная' THEN 2 
					WHEN r.type = 'трёхместная' THEN 3 
					WHEN r.type = 'двухместная (комфорт)' THEN 2 
					WHEN r.type = 'трёхместная (комфорт)' THEN 3 
					ELSE 0 
				END) AS user_count,
			h.number
		FROM Rooms r
		JOIN Hostels h ON r.Hostels_id = h.id
		LEFT JOIN Users u ON u.Rooms_id = r.id
		WHERE r.id = ?
		GROUP BY r.id, r.type, r.status, r.number, h.number`

	row := r.db.QueryRow(query, roomID)

	room := models.Room{}
	err := row.Scan(&room.ID, &room.Type, &room.Status, &room.Number, &room.UserCount, &room.HostelNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Room{}, fmt.Errorf("комната не найдена")
		}
		return models.Room{}, err
	}

	return room, nil
}

func (r *roomRepository) GetRoomNumberByID(roomID int) (int, error) {
	var roomNumber int
	query := "SELECT number FROM Rooms WHERE id = ?"
	err := r.db.QueryRow(query, roomID).Scan(&roomNumber)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении номера комнаты: %w", err)
	}
	return roomNumber, nil
}

func (r *roomRepository) GetRoomByNumber(roomNumber int) (*models.Room, error) {
	query := "SELECT * FROM Rooms WHERE number = ?"
	row := r.db.QueryRow(query, roomNumber)

	room := &models.Room{}
	err := row.Scan(&room.ID, &room.Type, &room.Status, &room.Number, &room.UserCount, &room.Hostel_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("комната не найдена")
		}
		return nil, err
	}

	return room, nil
}

func (r *roomRepository) GetRoomID(roomNumber string) (int, error) {
	var roomID int
	query := "SELECT id FROM Rooms WHERE number = ?"
	err := r.db.QueryRow(query, roomNumber).Scan(&roomID)
	if err != nil {
		return 0, errors.New("комната не найдена")
	}
	return roomID, nil
}

func (r *roomRepository) GetResidentsByRoomID(roomID int) ([]models.User, error) {
	query := `SELECT id, name, surname, email, institute, role FROM Users WHERE Rooms_id = ?`
	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	residents := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Surname, &user.Email, &user.Institute, &user.Role)
		if err != nil {
			return nil, err
		}
		residents = append(residents, user)
	}

	return residents, nil
}

func (r *roomRepository) InsertResidentIntoRoom(roomID int, email string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldRoomID *int
	query := "SELECT Rooms_id FROM Users WHERE email = ?;"
	err = tx.QueryRow(query, email).Scan(&oldRoomID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	updateUserQuery := "UPDATE Users SET Rooms_id = ? WHERE email = ?;"
	_, err = tx.Exec(updateUserQuery, roomID, email)
	if err != nil {
		return err
	}

	incrementRoomQuery := "UPDATE Rooms SET user_count = user_count + 1 WHERE id = ?;"
	_, err = tx.Exec(incrementRoomQuery, roomID)
	if err != nil {
		return err
	}

	if oldRoomID != nil {
		decrementRoomQuery := "UPDATE Rooms SET user_count = user_count - 1 WHERE id = ? AND user_count > 0;"
		_, err = tx.Exec(decrementRoomQuery, *oldRoomID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *roomRepository) DeleteResidentFromRoom(email string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}

	var currentRoomID int
	err = tx.QueryRow("SELECT Rooms_id FROM Users WHERE email = ?", email).Scan(&currentRoomID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при поиске пользователя по email '%s': %v", email, err)
	}

	if currentRoomID != 0 {
		_, err = tx.Exec("UPDATE Rooms SET user_count = user_count - 1 WHERE id = ?", currentRoomID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка при обновлении количества жильцов в комнате с id %d: %v", currentRoomID, err)
		}
	}

	_, err = tx.Exec("UPDATE Users SET Rooms_id = 999 WHERE email = ?", email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при перемещении пользователя с email '%s' в фантомную комнату: %v", email, err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при коммите транзакции: %v", err)
	}

	return nil
}

func (r *roomRepository) GetRoomIdByEmail(email string) (int, error) {
	var roomID int
	err := r.db.QueryRow("SELECT Rooms_id FROM Users WHERE email = ?", email).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func (r *roomRepository) GetRoomsCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM Rooms").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *roomRepository) GetHostelIDByRoomID(roomID int) (int, error) {
	var hostelNumber int
	query := "SELECT Hostels_id FROM Rooms WHERE id = ?"
	err := r.db.QueryRow(query, roomID).Scan(&hostelNumber)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении номера общежития: %w", err)
	}
	return hostelNumber, nil
}

func (r *roomRepository) GetHostelNumberByID(hostelID int) (int, error) {
	var hostelNumber int
	query := "SELECT number FROM Hostels WHERE id = ?"
	err := r.db.QueryRow(query, hostelID).Scan(&hostelNumber)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении номера общежития: %w", err)
	}
	return hostelNumber, nil
}

func (r *roomRepository) UpdateRoomStatus(roomNumber int, status string) error {
	_, err := r.db.Exec("UPDATE Rooms SET status = ? WHERE number = ?", status, roomNumber)
	return err
}

func (r *roomRepository) FreezeRoom(roomID int) error {
	_, err := r.db.Exec("UPDATE Rooms SET status = 'renovation' WHERE id = ? AND (SELECT COUNT(*) FROM Users WHERE room_id = Rooms.id) = 0", roomID)
	return err
}

func (r *roomRepository) GetInventoryByRoomID(id int) ([]models.Inventory, error) {
	return r.inventoryRepo.GetInventoryByRoomID(id)
}

func (r *roomRepository) GetRoomNumberByRoomID(roomID int) (int, error) {
	log.Println(roomID)
	var roomNumber int
	err := r.db.QueryRow("SELECT r.number FROM Rooms r LEFT JOIN Users u ON r.id = u.Rooms_id WHERE u.Rooms_id = ?", roomID).Scan(&roomNumber)
	if err != nil {
		return 0, err
	}
	return roomNumber, nil
}
