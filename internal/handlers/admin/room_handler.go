package handlers

import (
	"fmt"
	"hostel-management/internal/services"
	"hostel-management/pkg/middlewares"
	handlers "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
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

	const op = "handlers.room_handler.RoomsHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	var rooms []models.Room
	if role == "admin" {
		rooms, err = h.roomService.GetAllRooms()
		if err != nil {
			log.Printf("Unable to fetch rooms: %v: %v", err, op)
			middlewares.HandleError(c, 500, "Ошибка: не удалось получить комнаты")
			return
		}
	} else if role == "headman" {
		rooms, err = h.roomService.GetAllRoomsByHeadman(email)
		if err != nil {
			log.Printf("Unable to fetch rooms: %v: %v", err, op)
			middlewares.HandleError(c, 500, "Ошибка: не удалось получить комнаты")
			return
		}
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "admin_rooms",
		"Role":    role,
		"Rooms":   rooms,
		"Flashes": flashes,
	})
}

func (h *RoomHandler) RoomInfoHandler(c *gin.Context) {

	const op = "handlers.room_handler.RoomInfoHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for room: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	room, err := h.roomService.GetRoomByID(idInt)
	if err != nil {
		log.Printf("Failed to get room: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить комнату")
		return
	}

	residents, err := h.roomService.GetResidentsByRoomID(idInt)
	if err != nil {
		log.Printf("Failed to get residents: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить жильцов")
		return
	}

	inventory, err := h.roomService.GetInventoryByRoomID(idInt)
	if err != nil {
		log.Printf("Failed to get furniture: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить мебель")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":      "room",
		"Role":      role,
		"Room":      room,
		"Residents": residents,
		"Inventory": inventory,
		"Flashes":   flashes,
	})
}

func (h *RoomHandler) AddRoomHandler(c *gin.Context) {
	const op = "handlers.room_handler.AddRoomHandler"
	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	number, err := strconv.Atoi(c.PostForm("roomNumber"))
	if err != nil {
		log.Printf("Invalid room number: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный номер комнаты")
		return
	}
	hostelNumber, err := strconv.Atoi(c.PostForm("roomHostel"))
	if err != nil {
		log.Printf("Invalid hostel number: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный номер общежития")
		return
	}
	roomType := c.PostForm("roomType")
	roomStatus := c.PostForm("roomStatus")

	err = h.roomService.CreateRoom(roomType, roomStatus, number, 0, hostelNumber)
	if err != nil {
		log.Printf("Failed to create room: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось создать комнату")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Комната успешно создана!")
	session.Save()

	c.Redirect(303, "/admin/rooms")
}

func (h *RoomHandler) AddResidentIntoRoomHandler(c *gin.Context) {

	const op = "handlers.room_handler.AddResidentIntoRoomHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	roomID := c.Param("id")
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID комнаты")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	email := c.PostForm("email")

	err = h.roomService.InsertResidentIntoRoom(roomIDInt, email)
	if err != nil {
		log.Printf("Failed to add resident into room: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось добавить жильца в комнату")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Пользователь успешно добавлен!")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/rooms/room_info/%d", role, roomIDInt))
}

func (h *RoomHandler) DeleteResidentFromRoomHandler(c *gin.Context) {
	const op = "handlers.room_handler.DeleteResidentFromRoomHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	email := c.PostForm("email")
	if email == "" {
		log.Printf("Email is empty: %v", op)
		middlewares.HandleError(c, 400, "Ошибка: email пустой")
		return
	}

	roomID, err := h.roomService.DeleteResidentFromRoom(email)
	if err != nil {
		log.Printf("Failed to delete resident from room: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить жильца из комнаты")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Пользователь успешно удален!")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/rooms/room_info/%d", role, roomID))
}

func (h *RoomHandler) FreezeRoomHandler(c *gin.Context) {
	const op = "handlers.room_handler.FreezeRoomHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		log.Printf("Invalid room ID: %v", op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID комнаты")
		return
	}
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID комнаты")
		return
	}
	err = h.roomService.FreezeRoom(roomIDInt)
	if err != nil {
		log.Printf("Failed to freeze room: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось заморозить комнату")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Комната не активна")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/rooms/room_info/%s", role, roomID))
}

func (h *RoomHandler) UnfreezeRoomHandler(c *gin.Context) {
	const op = "handlers.room_handler.UnfreezeRoomHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		log.Printf("Invalid room ID: %v", op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID комнаты")
		return
	}
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID комнаты")
		return
	}
	err = h.roomService.UnfreezeRoom(roomIDInt)
	if err != nil {
		log.Printf("Failed to unfreeze room: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось разморозить комнату")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Комната активна")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/rooms/room_info/%s", role, roomID))
}
