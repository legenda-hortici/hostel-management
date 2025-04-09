package handlers

import (
	"fmt"
	"hostel-management/pkg/validation"
	"hostel-management/internal/services"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService services.RoomService
}

func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (h *RoomHandler) RoomsHandler(c *gin.Context) {
	const op = "handlers.RoomHandler.RoomsHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	rooms, err := h.roomService.GetAllRooms()
	if err != nil {
		log.Printf("Unable to fetch rooms: %v: %v", err, op)
		c.String(500, "Unable to fetch rooms")
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":  "admin_rooms",
		"Role":  "admin",
		"Rooms": rooms,
	})
}

func (h *RoomHandler) RoomInfoHandler(c *gin.Context) {

	const op = "handlers.RoomHandler.RoomInfoHandler"

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for room: %v: %v", err, op)
		c.String(400, "ID не найден в URL")
		return
	}

	room, err := h.roomService.GetRoomByID(idInt)
	if err != nil {
		log.Printf("Failed to get room: %v: %v", err, op)
		c.String(500, "Failed to get room")
		return
	}

	residents, err := h.roomService.GetResidentsByRoomID(idInt)
	if err != nil {
		log.Printf("Failed to get residents: %v: %v", err, op)
		c.String(500, "Failed to get residents")
		return
	}

	inventory, err := h.roomService.GetInventoryByRoomID(idInt)
	if err != nil {
		log.Printf("Failed to get furniture: %v: %v", err, op)
		c.String(500, "Failed to get furniture")
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":      "room",
		"Role":      "admin",
		"Room":      room,
		"Residents": residents,
		"Inventory": inventory,
	})
}

func (h *RoomHandler) AddRoomHandler(c *gin.Context) {
	const op = "handlers.RoomHandler.AddRoomHandler"
	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	number, err := strconv.Atoi(c.PostForm("roomNumber"))
	if err != nil {
		log.Printf("Invalid room number: %v: %v", err, op)
		c.String(400, "Invalid room number")
		return
	}
	hostelNumber, err := strconv.Atoi(c.PostForm("roomHostel"))
	if err != nil {
		log.Printf("Invalid hostel number: %v: %v", err, op)
		c.String(400, "Invalid hostel number")
		return
	}
	roomType := c.PostForm("roomType")
	roomStatus := c.PostForm("roomStatus")

	err = h.roomService.CreateRoom(roomType, roomStatus, number, 0, hostelNumber)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/admin/rooms")
}

func (h *RoomHandler) AddResidentIntoRoomHandler(c *gin.Context) {
	const op = "handlers.RoomHandler.AddResidentIntoRoomHandler"
	roomID := c.Param("id")
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v: %v", err, op)
		c.String(400, "Invalid room ID")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	email := c.PostForm("email")

	err = h.roomService.InsertResidentIntoRoom(roomIDInt, email)
	if err != nil {
		log.Printf("Failed to add resident into room: %v: %v", err, op)
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, fmt.Sprintf("/admin/rooms/room_info/%d", roomIDInt))
}

func (h *RoomHandler) DeleteResidentFromRoomHandler(c *gin.Context) {
	const op = "handlers.RoomHandler.DeleteResidentFromRoomHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	email := c.PostForm("email")
	if email == "" {
		log.Printf("Email is empty: %v", op)
		c.String(400, "Email is empty")
		return
	}

	roomID, err := h.roomService.DeleteResidentFromRoom(email)
	if err != nil {
		log.Printf("Failed to delete resident from room: %v: %v", err, op)
		c.String(500, err.Error())
		return
	}

	c.Redirect(303, fmt.Sprintf("/admin/rooms/room_info/%d", roomID))
}

func (h *RoomHandler) FreezeRoomHandler(c *gin.Context) {
	const op = "handlers.RoomHandler.FreezeRoomHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		log.Printf("Invalid room ID: %v", op)
		c.String(400, "Invalid room ID")
		return
	}
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v: %v", err, op)
		c.String(400, "Invalid room ID")
		return
	}
	err = h.roomService.FreezeRoom(roomIDInt)
	if err != nil {
		log.Printf("Failed to freeze room: %v: %v", err, op)
		c.String(500, "Failed to freeze room")
		return
	}

	c.Redirect(303, fmt.Sprintf("/room_info/%s", roomID))
}
