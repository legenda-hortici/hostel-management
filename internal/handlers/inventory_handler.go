package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"hostel-management/storage/models"
	"strconv"

	"log"

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

	const op = "handlers.InventoryHandler.InventoryHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v", op)
		return
	}

	inventory, err := h.inventoryService.GetAllInventory()
	if err != nil {
		c.String(500, err.Error())
		log.Printf("Failed to get inventory: %v: %v", err, op)
		return
	}

	for i := range inventory {
		inventory[i].Point = i + 1
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":      "admin_inventory",
		"Role":      role,
		"Inventory": inventory,
	})
}

func (h *InventoryHandler) DeleteInventoryItemHandler(c *gin.Context) {

	const op = "handlers.InventoryHandler.DeleteInventoryItemHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		log.Printf("Failed to get ID for inventory item: %v: %v", err, op)
		return
	}

	err = h.inventoryService.DeleteInventoryItem(id)
	if err != nil {
		c.String(500, "Failed to delete inventory item")
		log.Printf("Failed to delete inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/inventory")
}

func (h *InventoryHandler) AddInventoryItemHandler(c *gin.Context) {

	const op = "handlers.InventoryHandler.AddInventoryItemHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	furniture := c.PostForm("furniture")
	invNumber, err := strconv.Atoi(c.PostForm("inv_number"))
	if err != nil {
		c.String(400, "Invalid amount")
		log.Printf("Failed to get inv_number to add inventory item: %v: %v", err, op)
		return
	}
	room, err := strconv.Atoi(c.PostForm("room"))
	if err != nil {
		c.String(400, "Invalid room")
		log.Printf("Failed to get room to add inventory item: %v: %v", err, op)
		return
	}

	inventory := models.Inventory{
		Name:       furniture,
		Count:      1,
		InvNumber:  invNumber,
		RoomNumber: room,
	}

	err = h.inventoryService.InsertIntoInventory(inventory)
	if err != nil {
		c.String(500, "Failed to create inventory item")
		log.Printf("Failed to create inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/inventory")
}
