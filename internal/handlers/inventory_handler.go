package handlers

import (
	"hostel-management/internal/services"
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

	role, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	inventory, err := h.inventoryService.GetAllInventory()
	if err != nil {
		c.String(500, err.Error())
		log.Printf("Failed to get inventory: %v: %v", err, op)
		return
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

	err = h.inventoryService.InsertIntoInventory(furniture, invNumber, room)
	if err != nil {
		c.String(500, "Failed to create inventory item")
		log.Printf("Failed to create inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/inventory")
}
