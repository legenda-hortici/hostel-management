package handlers

import (
	"hostel-management/internal/services"
	handlers "hostel-management/pkg/validation"
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

	const op = "handlers.inventory_handler.InventoryHandler"

	role, err := handlers.ValidateUserByRole(c, op)
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

	const op = "handlers.inventory_handler.DeleteInventoryItemHandler"

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

	const op = "handlers.inventory_handler.AddInventoryItemHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	furniture := c.PostForm("furniture")
	invNumber := c.PostForm("inv_number")
	room, err := strconv.Atoi(c.PostForm("room"))
	if err != nil {
		c.String(400, "Invalid room")
		log.Printf("Failed to get room to add inventory item: %v: %v", err, op)
		return
	}
	hostel, err := strconv.Atoi(c.PostForm("hostel"))
	if err != nil {
		c.String(400, "Invalid hostel")
		log.Printf("Failed to get hostel to add inventory item: %v: %v", err, op)
		return
	}

	err = h.inventoryService.InsertIntoInventory(furniture, invNumber, room, hostel)
	if err != nil {
		c.String(500, "Failed to create inventory item")
		log.Printf("Failed to create inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/inventory")
}

func (h *InventoryHandler) UpdateInventoryItemHandler(c *gin.Context) {

	const op = "handlers.inventory_handler.UpdateInventoryItemHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.PostForm("id")
	if idStr == "" {
		c.String(400, "ID is required")
		log.Printf("ID is missing: %v", op)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		log.Printf("Failed to get ID for inventory item: %v: %v", err, op)
		return
	}

	furniture := c.PostForm("name")
	invNumber := c.PostForm("invnumber")
	room, err := strconv.Atoi(c.PostForm("roomnumber"))
	if err != nil {
		c.String(400, "Invalid room")
		log.Printf("Failed to get room to update inventory item: %v: %v", err, op)
		return
	}
	hostel, err := strconv.Atoi(c.PostForm("hostelnumber"))
	if err != nil {
		c.String(400, "Invalid hostel")
		log.Printf("Failed to get hostel to update inventory item: %v: %v", err, op)
		return
	}

	err = h.inventoryService.UpdateInventoryItem(id, furniture, invNumber, room, hostel)
	if err != nil {
		c.String(500, "Failed to update inventory item")
		log.Printf("Failed to update inventory item: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/inventory")
}
