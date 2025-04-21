package handlers

import (
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/internal/services"
	"hostel-management/pkg/helpers"
	"hostel-management/pkg/middlewares"
	validators "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
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
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	adminData, err := h.userService.GetAdminData(role)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить данные администратора")
		return
	}

	hostelData, err := h.hostelService.GetHostelsInfo(db.DB)
	if err != nil {
		log.Printf("Failed to get hostels info: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить информацию об общежитиях")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
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
		"Flashes": flashes,
	})
}

func (h *AdminHandler) UpdateCabinetHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.UpdateCabinetHandler"

	_, err := validators.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
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
		middlewares.HandleError(c, 500, "Ошибка: не удалось обновить данные администратора")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Данные администратора успешно обновлены!")
	session.Save()

	c.Redirect(303, "/admin")
}

func (h *AdminHandler) HostelInfoHandler(c *gin.Context) {

	const op = "handlers.AdminHandler.HostelInfoHandler"

	_, err := validators.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	hostelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get hostel ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	hostelInfo, err := h.hostelService.GetHostelInfo(hostelID)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить информацию о хостеле")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "hostel_info",
		"Role":    "admin",
		"Hostel":  hostelInfo,
		"Flashes": flashes,
	})
}

func (h *AdminHandler) AssignCommandantHandler(c *gin.Context) {

	const op = "handlers.admin.AssignCommandantHandler"

	_, err := validators.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	hostelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get hostel ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	email := c.PostForm("email")

	err = h.hostelService.InsertHeadmanIntoHostel(hostelID, email)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось назначить команданта")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Командант успешно назначен!")
	session.Save()

	c.Redirect(303, "/admin/hostel/"+fmt.Sprint(hostelID))
}

func (h *AdminHandler) RemoveCommandantHandler(c *gin.Context) {

	const op = "handlers.admin.RemoveCommandantHandler"

	_, err := validators.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	hostelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get hostel ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	err = h.hostelService.DeleteHeadmanFromHostel(hostelID)
	if err != nil {
		log.Printf("Failed to delete headman from hostel: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить команданта")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Командант успешно удален!")
	session.Save()

	c.Redirect(303, "/admin/hostel/"+fmt.Sprint(hostelID))
}

func CreateContractHandler(c *gin.Context) {

	const op = "handlers.admin.CreateContractHandler"

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 303, "Ошибка: метод не разрешен")
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
		middlewares.HandleError(c, 500, "Ошибка: не удалось создать договор")
		log.Printf("Failed to create contract: %v: %v", err, op)
		return
	}

	// Отправляем файл в ответе
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", "attachment; filename=Договор.docx")
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	c.Data(201, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", fileBytes)
}
