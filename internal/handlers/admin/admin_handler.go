package handlers

import (
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/internal/services"
	handlers "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userService   services.UserService
	hostelService services.HostelService
}

func NewAdminHandler(userService services.UserService, hostelService services.HostelService) *AdminHandler {
	return &AdminHandler{
		userService:   userService,
		hostelService: hostelService,
	}
}

func (h *AdminHandler) AdminCabinetHandler(c *gin.Context) {

	const op = "handlers.admin.AdminHandler.AdminCabinetHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	adminData, err := h.userService.GetAdminData(role)
	if err != nil {
		c.String(500, err.Error()+": "+op)
		return
	}

	hostelData, err := h.hostelService.GetHostelsInfo(db.DB)
	if err != nil {
		log.Printf("Failed to get hostels info: %v: %v", err, op)
		c.String(500, err.Error()+": "+op)
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page": "admin_cabinet",
		"Role": role,
		"Admin": map[string]interface{}{
			"Username": adminData.Username,
			"Surname":  adminData.Surname,
			"Email":    adminData.Email,
			"Role":     adminData.Role,
		},
		"Hostels": hostelData,
	})
}

func (h *AdminHandler) UpdateCabinetHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.UpdateCabinetHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	username := c.PostForm("username")
	surname := c.PostForm("surname")
	password := c.PostForm("password")

	var req models.UserRequest

	req.Username = username
	req.Surname = surname
	req.Password = password

	err = h.userService.UpdateAdminData(req)
	if err != nil {
		c.String(500, err.Error()+": "+op)
		return
	}

	c.Redirect(303, "/admin")
}

func (h *AdminHandler) HostelInfoHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.HostelInfoHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	hostelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get hostel ID: %v: %v", err, op)
		c.String(400, "Invalid hostel ID")
		return
	}

	fmt.Println(hostelID)

	// hostel, err := h.hostelService.GetHostelInfo(hostelID)
	// if err != nil {
	// 	c.String(500, err.Error()+": "+op)
	// 	return
	// }

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page": "hostel_info",
		"Role": "admin",
		// "Hostel": hostel,
	})
}
