package handlers

import (
	"fmt"
	"hostel-management/internal/services"
	"hostel-management/pkg/middlewares"
	handlers "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"strconv"

	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService services.InventoryService
}

func NewInventoryHandler(inventoryService services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) InventoryHandler(c *gin.Context) {

	const op = "handlers.inventory_handler.InventoryHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("acces denied: %v", err)
		middlewares.HandleError(c, 303, "Ошибка: доступ запрещен")
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("acces denied: %v", err)
		middlewares.HandleError(c, 303, "Ошибка: доступ запрещен")
		return
	}

	var inventory []models.Inventory
	if role == "admin" {
		inventory, err = h.inventoryService.GetAllInventory()
		if err != nil {
			middlewares.HandleError(c, 303, "Ошибка: не удалось получить инвентарь")
			log.Printf("Failed to get inventory: %v: %v", err, op)
			return
		}
	} else if role == "headman" {
		inventory, err = h.inventoryService.GetAllInventoryByHeadman(email)
		if err != nil {
			middlewares.HandleError(c, 303, "Ошибка: не удалось получить инвентарь")
			log.Printf("Failed to get inventory: %v: %v", err, op)
			return
		}
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":      "admin_inventory",
		"Role":      role,
		"Inventory": inventory,
		"Flashes":   flashes,
	})
}

func (h *InventoryHandler) DeleteInventoryItemHandler(c *gin.Context) {

	const op = "handlers.inventory_handler.DeleteInventoryItemHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("acces denied: %v", err)
		middlewares.HandleError(c, 303, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 303, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: ID не найден в URL")
		log.Printf("Failed to get ID for inventory item: %v: %v", err, op)
		return
	}

	err = h.inventoryService.DeleteInventoryItem(id)
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: не удалось удалить инвентарь")
		log.Printf("Failed to delete inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, fmt.Sprintf("/%s/inventory", role))
}

func (h *InventoryHandler) AddInventoryItemHandler(c *gin.Context) {

	const op = "handlers.inventory_handler.AddInventoryItemHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: доступ запрещен")
		log.Printf("acces denied: %v", err)
		return
	}

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 303, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	furniture := c.PostForm("furniture")
	invNumber := c.PostForm("inv_number")
	room, err := strconv.Atoi(c.PostForm("room"))
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: Комната не выбрана")
		log.Printf("Failed to get room to add inventory item: %v: %v", err, op)
		return
	}

	hostel, err := strconv.Atoi(c.PostForm("hostel"))
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: Общежитие не выбрано")
		log.Printf("Failed to get hostel to add inventory item: %v: %v", err, op)
		return
	}

	err = h.inventoryService.InsertIntoInventory(furniture, invNumber, room, hostel)
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: "+err.Error())
		log.Printf("Failed to add inventory item: %v: %v", err, op)
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Инвентарь успешно добавлен!")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/inventory", role))
}

func (h *InventoryHandler) UpdateInventoryItemHandler(c *gin.Context) {

	const op = "handlers.inventory_handler.UpdateInventoryItemHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("acces denied: %v", err)
		middlewares.HandleError(c, 303, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 303, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.PostForm("id")
	if idStr == "" {
		middlewares.HandleError(c, 303, "Ошибка: ID не найден")
		log.Printf("ID is missing: %v", op)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: ID не найден в URL")
		log.Printf("Failed to get ID for inventory item: %v: %v", err, op)
		return
	}

	furniture := c.PostForm("name")
	invNumber := c.PostForm("invnumber")
	room, err := strconv.Atoi(c.PostForm("roomnumber"))
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: Комната не выбрана")
		log.Printf("Failed to get room to update inventory item: %v: %v", err, op)
		return
	}

	hostel, err := strconv.Atoi(c.PostForm("hostelnumber"))
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: Общежитие не выбрано")
		log.Printf("Failed to get hostel to update inventory item: %v: %v", err, op)
		return
	}

	err = h.inventoryService.UpdateInventoryItem(id, furniture, invNumber, room, hostel)
	if err != nil {
		middlewares.HandleError(c, 303, "Ошибка: не удалось обновить инвентарь")
		log.Printf("Failed to update inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, fmt.Sprintf("/%s/inventory", role))
}
