package handlers

import (
	"fmt"
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/storage/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService services.RoomService
	roomHelper  helpers.RoomHelper
}

func NewRoomHandler(roomService services.RoomService, roomHelper helpers.RoomHelper) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		roomHelper:  roomHelper,
	}
}

func (h *RoomHandler) RoomsHandler(c *gin.Context) {
	rooms, err := h.roomService.GetAllRooms()
	if err != nil {
		c.String(500, "RoomsHandler: Unable to fetch rooms")
		return
	}

	for room := range rooms {
		if rooms[room].Status != "renovation" {
			h.roomHelper.DefineRoomStatus(rooms[room].Type, rooms[room].UserCount, rooms[room].Number)
		}
		rooms[room] = h.roomHelper.TranslateRoom(rooms[room])
	}

	data := map[string]interface{}{
		"Page":  "admin_rooms",
		"Role":  "admin",
		"Rooms": rooms,
	}

	c.HTML(200, "layout.html", data)
}

func (h *RoomHandler) RoomInfoHandler(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "ID не найден в URL")
		return
	}

	room, err := h.roomService.GetRoomByID(idInt)
	if err != nil {
		c.String(500, "RoomHandler: Failed to get room")
		return
	}

	residents, err := h.roomService.GetResidentsByRoomID(idInt)
	if err != nil {
		c.String(500, "RoomHandler: Failed to get residents")
		return
	}
	room.UserCount = len(residents)

	room = h.roomHelper.TranslateRoom(room)

	inventory, err := h.roomService.GetInventoryByRoomID(idInt)
	if err != nil {
		c.String(500, "RoomHandler: Failed to get furniture")
		return
	}

	data := map[string]interface{}{
		"Page":      "room",
		"Role":      "admin",
		"Room":      room,
		"Residents": residents,
		"Inventory": inventory,
	}
	c.HTML(200, "layout.html", data)
}

func (h *RoomHandler) AddRoomHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	number, err := strconv.Atoi(c.PostForm("roomNumber"))
	if err != nil {
		c.String(400, "Invalid room number")
		return
	}

	hostelNumber, err := strconv.Atoi(c.PostForm("roomHostel"))
	if err != nil {
		c.String(400, "Invalid hostel number")
		return
	}

	validTypes := map[string]bool{
		"once":           true,
		"double":         true,
		"triple":         true,
		"premium double": true,
		"premium triple": true,
	}

	validStatuses := map[string]bool{
		"unoccupied": true,
		"occupied":   true,
		"renovation": true,
	}

	roomType := c.PostForm("roomType")
	if !validTypes[roomType] {
		c.String(400, "Invalid room type")
		return
	}

	roomStatus := c.PostForm("roomStatus")
	if !validStatuses[roomStatus] {
		c.String(400, "Invalid room status")
		return
	}

	room := models.Room{
		Number:       number,
		Type:         roomType,
		Status:       roomStatus,
		HostelNumber: hostelNumber,
		UserCount:    0,
	}

	err = h.roomHelper.ValidateRoomData(room)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/admin/rooms")
}

func (h *RoomHandler) AddResidentIntoRoomHandler(c *gin.Context) {
	roomID := c.Param("id")
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		c.String(400, "Invalid room ID")
		return
	}

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	email := c.PostForm("email")

	err = h.roomHelper.ValidateAddResidentData(email, roomIDInt)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, fmt.Sprintf("/admin/rooms/room_info/%d", roomIDInt))
}

func (h *RoomHandler) DeleteResidentFromRoomHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	email := c.PostForm("email")
	if email == "" {
		c.String(400, "Email is empty")
		return
	}

	roomID, err := h.roomHelper.ValidateDeleteResidentData(email)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(303, fmt.Sprintf("/admin/rooms/room_info/%d", roomID))
}

func (h *RoomHandler) FreezeRoomHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		c.String(400, "FreezeRoomHandler: ID не найден в URL")
		return
	}

	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		c.String(400, "FreezeRoomHandler: Некорректный ID")
		return
	}

	err = h.roomService.FreezeRoom(roomIDInt)
	if err != nil {
		c.String(500, "FreezeRoomHandler: Невозможно заморозить комнату, так как в ней есть жильцы")
		return
	}

	c.Redirect(303, fmt.Sprintf("/room_info/%d", roomIDInt))
}
