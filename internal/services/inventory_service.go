package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
	"log"
)

type InventoryService interface {
	GetAllInventory() ([]models.Inventory, error)
	InsertIntoInventory(inventory models.Inventory) error
	DeleteInventoryItem(id int) error
	GetInventoryByRoomID(roomID int) ([]models.Inventory, error)
}

type inventoryService struct {
	inventoryRepo repositories.InventoryRepository
}

func NewInventoryService(inventoryRepo repositories.InventoryRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: repositories.NewInventoryRepository(),
	}
}

func (s *inventoryService) GetAllInventory() ([]models.Inventory, error) {
	return s.inventoryRepo.GetAllInventory()
}

func (s *inventoryService) InsertIntoInventory(inventory models.Inventory) error {
	log.Println(inventory.Icon)
	switch inventory.Name {
	case "Стул":
		inventory.Icon = "img/svg/chair-svgrepo-com.svg"
	case "Стол":
		inventory.Icon = "img/svg/desk-svgrepo-com.svg"
	case "Кровать":
		inventory.Icon = "img/svg/bed-svgrepo-com.svg"
	case "Тумбочка":
		inventory.Icon = "img/svg/сhest-of-drawers-svgrepo-com.svg"
	case "Шкаф":
		inventory.Icon = "img/svg/wardrobe-svgrepo-com.svg"
	case "Стеллаж":
		inventory.Icon = "img/svg/bookshelf-svgrepo-com.svg"
	}
	return s.inventoryRepo.InsertIntoInventory(inventory)
}

func (s *inventoryService) DeleteInventoryItem(id int) error {
	return s.inventoryRepo.DeleteInventory(id)
}

func (s *inventoryService) GetInventoryByRoomID(roomID int) ([]models.Inventory, error) {
	return s.inventoryRepo.GetInventoryByRoomID(roomID)
}
