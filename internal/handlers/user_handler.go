package handlers

import (
	"database/sql"
	"fmt"
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/internal/session"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
	roomService services.RoomService
	userHelper  helpers.UserHelper
}

func NewUserHandler(userService services.UserService, roomService services.RoomService, userHelper helpers.UserHelper) *UserHandler {
	return &UserHandler{
		userService: userService,
		roomService: roomService,
		userHelper:  userHelper,
	}
}

func (h *UserHandler) AdminCabinetHandler(c *gin.Context) {
	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied")
		return
	}

	adminData, err := h.userService.GetAdminData(role)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	log.Println(adminData)
	data := map[string]interface{}{
		"Page":  "admin_cabinet",
		"Role":  role,
		"Admin": adminData,
	}

	c.HTML(200, "layout.html", data)
}

func (h *UserHandler) ResidentsHandler(c *gin.Context) {
	searchTerm := c.Query("search")

	residents, err := h.userService.GetAllUsers()
	if err != nil {
		c.String(500, "ResidentsHandler: Failed to get residents")
		return
	}

	for i := range residents {
		if residents[i].Room_id > 0 {
			room, err := h.roomService.GetRoomByID(residents[i].Room_id)
			if err != nil {
				c.String(500, "ResidentsHandler: Failed to get room number")
				return
			}
			residents[i].RoomNumber = room.Number
		} else {
			residents[i].RoomNumber = 0
		}
	}

	data := map[string]interface{}{
		"Page":       "admin_residents",
		"Role":       "admin",
		"Residents":  residents,
		"SearchTerm": searchTerm,
	}

	c.HTML(200, "layout.html", data)
}

func (h *UserHandler) ResidentInfoHandler(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "ID не найден в URL")
		return
	}

	resident, err := h.userService.GetUserByID(idInt)
	if err != nil {
		c.String(500, "ResidentsHandler: Failed to get resident")
		return
	}
	resident.RoomNumber, _ = h.roomService.GetRoomNumberByID(resident.Room_id)

	data := map[string]interface{}{
		"Page":     "resident",
		"Role":     "admin",
		"Resident": resident,
	}
	c.HTML(200, "layout.html", data)
}

func (h *UserHandler) AddResidentHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "AddResident: Method not allowed")
		return
	}

	institute := c.PostForm("institute")
	roomNumber := c.PostForm("room")
	roomNumberInt, err := strconv.Atoi(roomNumber)
	if err != nil {
		c.String(400, "Некорректный номер комнаты")
		return
	}

	resident := models.User{
		Username:   c.PostForm("username"),
		Email:      c.PostForm("email"),
		Password:   c.PostForm("password"),
		Institute:  sql.NullString{String: institute, Valid: institute != ""},
		RoomNumber: roomNumberInt,
	}

	err = h.userHelper.ValidateResidentData(resident)
	if err != nil {
		c.String(400, "Ошибка добавления жильца: "+err.Error())
		return
	}

	c.Redirect(303, "/admin/residents")
}

func (h *UserHandler) UpdateResidentDataHandler(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.String(405, "UpdateResidentData: Method not allowed")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Некорректный ID")
		return
	}

	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Institute string `json:"institute"`
		Role      string `json:"role"`
		Password  string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(400, "Ошибка парсинга JSON")
		return
	}

	resident := models.User{
		Username: req.Username,
		Email:    req.Email,
		Institute: sql.NullString{
			String: req.Institute,
			Valid:  req.Institute != "",
		},
		Role:     req.Role,
		Password: req.Password,
	}

	err = h.userService.UpdateUser(&resident)
	if err != nil {
		c.String(500, "UpdateResidentData: Failed to update resident data")
		return
	}

	c.Redirect(303, fmt.Sprintf("/admin/residents/%d", idInt))
}

func (h *UserHandler) DeleteResidentHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "DeleteResident: Method not allowed")
		return
	}

	if c.PostForm("_method") != "DELETE" {
		c.String(400, "Неверный метод")
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Некорректный ID")
		return
	}

	err = h.userService.DeleteUser(idInt)
	if err != nil {
		c.String(500, "DeleteResident: Failed to delete resident")
		return
	}

	c.Redirect(303, "/admin/residents")
}
