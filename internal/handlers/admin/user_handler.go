package handlers

import (
	"fmt"
	"hostel-management/internal/services"
	"hostel-management/pkg/middlewares"
	validate "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
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

	role, err := validate.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	email, err := validate.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	searchTerm := c.Query("search")

	var residents []models.User

	if role == "admin" {
		residents, err = h.userService.GetAllUsers()
		if err != nil {
			log.Printf("Failed to get residents: %v: %v", err, op)
			c.String(500, "Failed to get residents")
			return
		}
	} else if role == "headman" {
		residents, err = h.userService.GetAllByHeadman(email)
		if err != nil {
			log.Printf("Failed to get residents: %v: %v", err, op)
			c.String(500, "Failed to get residents")
			return
		}
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":       "admin_residents",
		"Role":       role,
		"Residents":  residents,
		"SearchTerm": searchTerm,
	})
}

func (h *UserHandler) ResidentInfoHandler(c *gin.Context) {

	const op = "handlers.ResidentInfoHandler.ResidentInfoHandler"

	role, err := validate.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for resident: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить ID жильца")
		return
	}

	resident, err := h.userService.GetUserByID(idInt)
	if err != nil {
		log.Printf("Failed to get resident: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить жильца")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":     "resident",
		"Role":     role,
		"Resident": resident,
		"Flashes":  flashes,
	})
}

func (h *UserHandler) AddResidentHandler(c *gin.Context) {

	const op = "handlers.user_handler.AddResidentHandler"

	role, err := validate.ValidateUserByRole(c, op)
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

	username := c.PostForm("username")
	surname := c.PostForm("surname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	institute := c.PostForm("institute")
	roomNumber := c.PostForm("room")
	roomNumberInt, err := strconv.Atoi(roomNumber)
	if err != nil {
		log.Printf("Invalid room number: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный номер комнаты")
		return
	}

	err = h.userService.CreateUser(username, surname, email, password, institute, roomNumberInt)
	if err != nil {
		log.Printf("Failed to add resident: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось добавить жильца")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Житель успешно добавлен!")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/residents", role))
}

func (h *UserHandler) UpdateResidentDataHandler(c *gin.Context) {

	const op = "handlers.UpdateResidentDataHandler.UpdateResidentDataHandler"

	if c.Request.Method != "PUT" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	var req models.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid JSON: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный JSON")
		return
	}

	err = h.userService.UpdateUser(idInt, &req)
	if err != nil {
		log.Printf("Failed to update resident data: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось обновить данные жильца")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Данные жильца успешно обновлены!")
	session.Save()

	c.JSON(200, gin.H{
		"message": "Resident data updated successfully",
		"status":  "success",
		"data":    req,
	})
}

func (h *UserHandler) DeleteResidentHandler(c *gin.Context) {

	const op = "handlers.DeleteResidentHandler.DeleteResidentHandler"

	role, err := validate.ValidateUserByRole(c, op)
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

	if c.PostForm("_method") != "DELETE" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	err = h.userService.DeleteUser(idInt)
	if err != nil {
		log.Printf("Failed to delete resident: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить жильца")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Житель успешно удален!")
	session.Save()

	c.Redirect(303, fmt.Sprintf("/%s/residents", role))
}
