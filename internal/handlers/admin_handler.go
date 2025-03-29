package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"hostel-management/storage/db"
	"log"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userService   services.UserService
	hostelService services.HostelService
}

func NewAdminHandler(userService services.UserService, hostelService services.HostelService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
		hostelService: hostelService,
	}
}

func (h *AdminHandler) AdminCabinetHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.AdminCabinetHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied: %v: %v", role, op)
		return
	}

	adminData, err := h.userService.GetAdminData(role)
	if err != nil {
		c.String(500, err.Error()+": "+op)
		return
	}

	hostels, err := h.hostelService.GetHostelsInfo(db.DB)
	if err != nil {
		log.Printf("Failed to get hostels info: %v: %v", err, op)
		c.String(500, err.Error()+": "+op)
		return
	}

	hostelData := []map[string]interface{}{}
	for _, hostel := range hostels {
		hostelData = append(hostelData, map[string]interface{}{
			"Number":         hostel.HostelNumber,
			"RoomsCount":     hostel.OccupiedRooms + hostel.AvailableRooms,
			"OccupiedRooms":  hostel.OccupiedRooms,
			"AvailableRooms": hostel.AvailableRooms,
			"ResidentsCount": hostel.ResidentsCount,
			"Location":       hostel.HostelLocation,
		})
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page": "admin_cabinet",
		"Role": role,
		"Admin": map[string]interface{}{
			"Username": adminData.Username,
			"Email":    adminData.Email,
			"Role":     adminData.Role,
		},
		"Hostels":  hostelData,
		"Messages": []map[string]interface{}{},
	})
}

func (h *AdminHandler) UpdateCabinetHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.UpdateCabinetHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied: %v: %v", role, op)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	err := h.userService.UpdateAdminData(username, password)
	if err != nil {
		c.String(500, err.Error()+": "+op)
		return
	}

	c.Redirect(303, "/admin")
}
