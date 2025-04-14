package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
)

type InventoryRepository interface {
	GetAllInventory() ([]models.Inventory, error)
	GetAllInventoryByHeadman(email string) ([]models.Inventory, error)
	InsertIntoInventory(inventory models.Inventory) error
	DeleteInventory(id int) error
	GetInventoryByRoomID(roomID int) ([]models.Inventory, error)
	UpdateInventoryItem(inventory models.Inventory) error
}

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository() InventoryRepository {
	return &inventoryRepository{
		db: db.DB,
	}
}

func (r *inventoryRepository) GetAllInventory() ([]models.Inventory, error) {
	query := `
		SELECT 
			i.id, 
			i.name, 
			i.inv_number, 
			r.number AS room_number, 
			i.icon, 
			h.number AS hostel_number
		FROM Inventory i
		JOIN Rooms r ON i.Rooms_id = r.id
		JOIN Hostels h ON r.Hostels_id = h.id;
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.Inventory
	for rows.Next() {
		var inv models.Inventory
		err := rows.Scan(
			&inv.ID,
			&inv.Name,
			&inv.InvNumber,
			&inv.RoomNumber,
			&inv.Icon,
			&inv.HostelNumber,
		)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, inv)
	}

	return inventory, nil
}

func (r *inventoryRepository) GetAllInventoryByHeadman(email string) ([]models.Inventory, error) {
	query := `
		SELECT 
			i.id, 
			i.name, 
			i.inv_number, 
			r.number AS room_number, 
			i.icon, 
			h.number AS hostel_number
		FROM 
			Inventory i
		JOIN 
			Rooms r ON i.Rooms_id = r.id
		JOIN 
			Hostels h ON r.Hostels_id = h.id
		WHERE 
			h.headman_email = ?;
	`

	rows, err := db.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.Inventory
	for rows.Next() {
		var inv models.Inventory
		err := rows.Scan(
			&inv.ID,
			&inv.Name,
			&inv.InvNumber,
			&inv.RoomNumber,
			&inv.Icon,
			&inv.HostelNumber,
		)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, inv)
	}

	return inventory, nil

}

func (r *inventoryRepository) InsertIntoInventory(inventory models.Inventory) error {
	query := `
		SELECT Rooms.id 
		FROM Rooms 
		JOIN Hostels ON Rooms.Hostels_id = Hostels.id 
		WHERE Rooms.number = ? AND Hostels.number = ?
	`

	err := db.DB.QueryRow(query, inventory.RoomNumber, inventory.HostelNumber).Scan(&inventory.Rooms_id)
	if err != nil {
		return fmt.Errorf("не найдена комната №%d в общежитии №%d: %w", inventory.RoomNumber, inventory.HostelNumber, err)
	}

	insertQuery := "INSERT INTO Inventory (name, inv_number, Rooms_id, icon) VALUES (?, ?, ?, ?)"
	_, err = db.DB.Exec(insertQuery, inventory.Name, inventory.InvNumber, inventory.Rooms_id, inventory.Icon)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении инвентаря: %w", err)
	}

	return nil
}

func (r *inventoryRepository) DeleteInventory(id int) error {
	query := "DELETE FROM Inventory WHERE id = ?"
	_, err := db.DB.Exec(query, id)
	return err
}

func (r *inventoryRepository) GetInventoryByRoomID(roomID int) ([]models.Inventory, error) {
	query := "SELECT * FROM Inventory WHERE Rooms_id = ?"
	rows, err := db.DB.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := []models.Inventory{}
	for rows.Next() {
		inv := models.Inventory{}
		err := rows.Scan(&inv.ID, &inv.Name, &inv.InvNumber, &inv.Rooms_id, &inv.Icon)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, inv)
	}

	return inventory, nil
}

func (r *inventoryRepository) UpdateInventoryItem(inventory models.Inventory) error {
	query := `
		SELECT Rooms.id 
		FROM Rooms 
		JOIN Hostels ON Rooms.Hostels_id = Hostels.id 
		WHERE Rooms.number = ? AND Hostels.number = ?
	`

	err := db.DB.QueryRow(query, inventory.RoomNumber, inventory.HostelNumber).Scan(&inventory.Rooms_id)
	if err != nil {
		return fmt.Errorf("не найдена комната №%d в общежитии №%d: %w", inventory.RoomNumber, inventory.HostelNumber, err)
	}

	query = "UPDATE Inventory SET name = ?, inv_number = ?, Rooms_id = ?, icon = ? WHERE id = ?"
	_, err = db.DB.Exec(query, inventory.Name, inventory.InvNumber, inventory.Rooms_id, inventory.Icon, inventory.ID)
	return err
}
