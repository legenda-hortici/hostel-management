package repositories

import (
	"database/sql"
	"hostel-management/storage/db"
	"hostel-management/storage/models"
)

type InventoryRepository interface {
	GetAllInventory() ([]models.Inventory, error)
	InsertIntoInventory(name string, count, invNumber, roomsID int) error
	DeleteInventory(id int) error
	GetInventoryByRoomID(roomID int) ([]models.Inventory, error)
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
	query := `SELECT i.id, i.name, i.count, i.inv_number, r.number 
				FROM Inventory i
				JOIN Rooms r ON i.Rooms_id = r.id;`
	var inventory []models.Inventory

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var inv models.Inventory
		err := rows.Scan(&inv.ID, &inv.Name, &inv.Count, &inv.InvNumber, &inv.RoomNumber)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, inv)
	}

	return inventory, nil
}

func (r *inventoryRepository) InsertIntoInventory(name string, count, invNumber, roomsID int) error {
	query := "INSERT INTO Inventory (name, count, inv_number, Rooms_id) VALUES (?, ?, ?, ?)"
	_, err := db.DB.Exec(query, name, count, invNumber, roomsID)
	return err
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
		err := rows.Scan(&inv.ID, &inv.Name, &inv.Count, &inv.InvNumber, &inv.Rooms_id)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, inv)
	}

	return inventory, nil
}
