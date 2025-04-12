package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
	roomService services.RoomService
}

func NewUserHandler(userService services.UserService, roomService services.RoomService) *UserHandler {
	return &UserHandler{
		userService: userService,
		roomService: roomService,
	}
}

func (h *UserHandler) ResidentsHandler(c *gin.Context) {

	const op = "handlers.ResidentsHandler.ResidentsHandler"

	searchTerm := c.Query("search")

	residents, err := h.userService.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get residents: %v: %v", err, op)
		c.String(500, "ResidentsHandler: Failed to get residents")
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":       "admin_residents",
		"Role":       "admin",
		"Residents":  residents,
		"SearchTerm": searchTerm,
	})
}

func (h *UserHandler) ResidentInfoHandler(c *gin.Context) {

	const op = "handlers.ResidentInfoHandler.ResidentInfoHandler"

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for resident: %v: %v", err, op)
		c.String(400, "ID not valid")
		return
	}

	resident, err := h.userService.GetUserByID(idInt)
	if err != nil {
		log.Printf("Failed to get resident: %v: %v", err, op)
		c.String(500, "Failed to get resident")
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":     "resident",
		"Role":     "admin",
		"Resident": resident,
	})
}

func (h *UserHandler) AddResidentHandler(c *gin.Context) {

	const op = "handlers.user_handler.AddResidentHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "AddResident: Method not allowed")
		return
	}

	username := c.PostForm("username")
	surname := c.PostForm("surname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	institute := c.PostForm("institute")
	roomNumber := c.PostForm("room")
	roomNumberInt, err := strconv.Atoi(roomNumber)
	if err != nil {
		log.Printf("Invalid room number: %v: %v", err, op)
		c.String(400, "Invalid room number")
		return
	}

	err = h.userService.CreateUser(username, surname, email, password, institute, roomNumberInt)
	if err != nil {
		log.Printf("Failed to add resident: %v: %v", err, op)
		c.String(400, "Failed to add resident: "+err.Error())
		return
	}

	c.Redirect(303, "/admin/residents")
}

func (h *UserHandler) UpdateResidentDataHandler(c *gin.Context) {

	const op = "handlers.UpdateResidentDataHandler.UpdateResidentDataHandler"

	if c.Request.Method != "PUT" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "UpdateResidentData: Method not allowed")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		c.String(400, "ID not valid")
		return
	}

	var req models.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid JSON: %v: %v", err, op)
		c.String(400, "Invalid JSON")
		return
	}

	err = h.userService.UpdateUser(idInt, &req)
	if err != nil {
		log.Printf("Failed to update resident data: %v: %v", err, op)
		c.String(500, "Failed to update resident data")
		return
	}

	c.JSON(200, gin.H{
		"message": "Resident data updated successfully",
		"status":  "success",
		"data":    req,
	})
}

func (h *UserHandler) DeleteResidentHandler(c *gin.Context) {

	const op = "handlers.DeleteResidentHandler.DeleteResidentHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	if c.PostForm("_method") != "DELETE" {
		log.Printf("Method not allowed: %v", op)
		c.String(400, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		c.String(400, "ID not valid")
		return
	}

	err = h.userService.DeleteUser(idInt)
	if err != nil {
		log.Printf("Failed to delete resident: %v: %v", err, op)
		c.String(500, "Failed to delete resident")
		return
	}

	c.Redirect(303, "/admin/residents")
}
