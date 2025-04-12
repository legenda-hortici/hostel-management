package handlers

import (
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/internal/services"
	"hostel-management/pkg/helpers"
	"hostel-management/pkg/session"
	validators "hostel-management/pkg/validation"
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

	role, err := validators.ValidateUserByRole(c, op)
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
			"Password": adminData.Password,
			"Email":    adminData.Email,
			"Role":     adminData.Role,
		},
		"Hostels": hostelData,
	})
}

func (h *AdminHandler) UpdateCabinetHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.UpdateCabinetHandler"

	_, err := validators.ValidateUserByRole(c, op)
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

	_, err := validators.ValidateUserByRole(c, op)
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

	hostelInfo, err := h.hostelService.GetHostelInfo(hostelID)
	if err != nil {
		c.String(500, err.Error()+": "+op)
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":   "hostel_info",
		"Role":   "admin",
		"Hostel": hostelInfo,
	})
}

func (h *AdminHandler) AssignCommandantHandler(c *gin.Context) {

	const op = "handlers.admin.AssignCommandantHandler"

	_, err := validators.ValidateUserByRole(c, op)
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

	email := c.PostForm("email")

	err = h.hostelService.InsertHeadmanIntoHostel(hostelID, email)
	if err != nil {
		c.String(500, err.Error()+": "+op)
		return
	}

	c.Redirect(303, "/admin/hostel/"+fmt.Sprint(hostelID))
}

func DocumentsHandler(c *gin.Context) {

	const op = "handlers.admin.DocumentsHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied: %v: %v", role, op)
		return
	}

	data := map[string]interface{}{
		"Page": "admin_documents",
		"Role": role,
	}

	c.HTML(200, "layout.html", data)
}

func CreateContractHandler(c *gin.Context) {

	const op = "handlers.admin.CreateContractHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed: %v: %v", c.Request.Method, op)
		return
	}

	data := helpers.ContractData{
		FirstName:    c.PostForm("firstName"),
		LastName:     c.PostForm("lastName"),
		MiddleName:   c.PostForm("middleName"),
		CheckInDate:  c.PostForm("checkInDate"),
		CheckOutDate: c.PostForm("checkOutDate"),
		RoomNumber:   c.PostForm("roomNumber"),
		Amount:       c.PostForm("amount"),
	}

	fileBytes, err := helpers.GenerateContract(data)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Отправляем файл в ответе
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", "attachment; filename=Договор.docx")
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	c.Data(201, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", fileBytes)
}
