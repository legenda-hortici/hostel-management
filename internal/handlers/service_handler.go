package handlers

import (
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceService   services.ServiceService
	userService      services.UserService
	statementService services.StatementService
	roomService      services.RoomService
}

func NewServiceHandler(serviceService services.ServiceService, userService services.UserService, statementService services.StatementService, roomService services.RoomService) *ServiceHandler {
	return &ServiceHandler{
		serviceService:   serviceService,
		userService:      userService,
		statementService: statementService,
		roomService:      roomService,
	}
}

func (h *ServiceHandler) ServicesHandler(c *gin.Context) {
	const op = "handlers.ServiceHandler.ServicesHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" && role != "user" {
		log.Printf("Access denied: %v", op)
		c.String(403, "Access denied")
		return
	}

	email, exists := session.GetUserEmail(c)
	if !exists || email == "" {
		log.Printf("User not authenticated: %v", op)
		c.String(401, "User not authenticated")
		return
	}

	services, err := h.serviceService.GetAllServices()
	if err != nil {
		c.String(500, "Ошибка получения сервисов: "+err.Error())
		return
	}

	for i := range services {
		services[i].Point = i + 1
	}

	statements, err := h.statementService.GetAllStatements()
	if err != nil {
		c.String(500, "Ошибка получения заявок: "+err.Error())
		return
	}

	userStatements, err := h.statementService.GetAllUserStatements(email)
	if err != nil {
		log.Printf("Failed to get user statements: %v: %v", err, op)
		c.String(500, "Ошибка получения заявок пользователя: "+err.Error())
		return
	}

	for statement := range statements {
		statements[statement].Status = helpers.TranslateStatus(statements[statement].Status)
		user, err := h.userService.GetUserByID(statements[statement].Users_id)
		if err != nil {
			c.String(500, "Ошибка получения пользователя: "+err.Error())
			return
		}
		statements[statement].Username = user.Username
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":           "services",
		"Role":           role,
		"Services":       services,
		"Statements":     statements,
		"userStatements": userStatements,
	})
}

func (h *ServiceHandler) AddServiceHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	name := c.PostForm("name")
	typeService := c.PostForm("type")
	description := c.PostForm("description")
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	isDate := c.PostForm("is_date") == "on"
	isHostel := c.PostForm("is_hostel") == "on"
	isPhone := c.PostForm("is_phone") == "on"

	err := h.serviceService.CreateService(name, typeService, description, isDate, isHostel, isPhone, amount)
	if err != nil {
		c.String(500, "Failed to create service")
		return
	}

	c.Redirect(303, "/admin/services")
}

func (h *ServiceHandler) ServiceInfoHandler(c *gin.Context) {
	role, exists := session.GetUserRole(c)
	if !exists {
		c.String(403, "Access denied")
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		return
	}

	service, err := h.serviceService.GetServiceByID(id)
	if err != nil {
		c.String(500, "Failed to get service")
		return
	}

	data := map[string]interface{}{
		"Page":    "service_info",
		"Role":    role,
		"Service": service,
	}
	c.HTML(200, "layout.html", data)
}

func (h *ServiceHandler) UpdateServiceHandler(c *gin.Context) {
	// Проверяем скрытое поле _method
	if c.Request.Method != "POST" || c.Request.FormValue("_method") != "PUT" {
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		return
	}

	service := models.Service{
		ID:          id,
		Name:        c.PostForm("name"),
		Type:        c.PostForm("type"),
		Description: c.PostForm("description"),
		Amount:      0,
		Is_date:     c.PostForm("is_date") == "on",
		Is_hostel:   c.PostForm("is_hostel") == "on",
		Is_phone:    c.PostForm("is_phone") == "on",
	}

	if amountStr := c.PostForm("amount"); amountStr != "" {
		service.Amount, err = strconv.Atoi(amountStr)
		if err != nil {
			c.String(400, "Invalid amount")
			return
		}
	}

	err = h.serviceService.UpdateServiceByID(id, service)
	if err != nil {
		c.String(500, "Failed to update service")
		return
	}

	c.Redirect(303, "/admin/services/service/"+idStr)
}

func (h *ServiceHandler) DeleteServiceHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		return
	}

	err = h.serviceService.DeleteService(id)
	if err != nil {
		c.String(500, "Failed to delete service")
		return
	}

	c.Redirect(303, "/admin/services")
}

func (h *ServiceHandler) RequestServiceHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	email, exists := session.GetUserEmail(c)
	if !exists || email == "" {
		c.String(401, "User not authenticated")
		return
	}

	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(400, "Invalid service ID")
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		c.String(500, "Failed to get user")
		return
	}

	amount, _ := strconv.Atoi(c.PostForm("amount"))
	request_date := c.PostForm("request_date")
	hostel, _ := strconv.Atoi(c.PostForm("hostel"))
	phone := c.PostForm("phone")

	statement := models.Statement{
		Name:     c.PostForm("name"),
		Type:     c.PostForm("type"),
		Amount:   amount,
		Date:     request_date,
		Phone:    phone,
		Status:   "awaits",
		Hostel:   hostel,
		Users_id: user.ID,
	}

	err = h.statementService.CreateStatementRequest(statement)
	if err != nil {
		c.String(500, "Failed to create service request")
		return
	}

	c.Redirect(303, "/services")
}

func (h *ServiceHandler) RequestInfoHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RequestInfoHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role == "" {
		log.Printf("Access denied: %v: %v", role, op)
		c.String(403, "Access denied: %v: %v", role, op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	request, err := h.statementService.GetStatementRequestByID(id)
	if err != nil {
		c.String(500, "Failed to get service request")
		return
	}

	user, err := h.userService.GetUserByID(request.Users_id)
	if err != nil {
		log.Printf("Failed to get user: %v: %v", err, op)
		c.String(500, "Failed to get user")
		return
	}

	roomNumber, err := h.roomService.GetRoomNumberByRoomID(user.Room_id)
	if err != nil {
		log.Printf("Failed to get room number: %v: %v", err, op)
		c.String(500, "Failed to get room number")
		return
	}
	user.RoomNumber = roomNumber

	data := map[string]interface{}{
		"Page":      "request_info",
		"Role":      role,
		"Statement": request,
		"User":      user,
	}
	c.HTML(200, "layout.html", data)
}

func (h *ServiceHandler) AcceptRequestHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		return
	}

	err = h.statementService.UpdateStatementRequestStatus(id, "approved")
	if err != nil {
		c.String(500, "Failed to accept request")
		return
	}

	c.Redirect(303, "/admin/services/request_info/"+idStr)
}

func (h *ServiceHandler) RejectRequestHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		return
	}

	err = h.statementService.UpdateStatementRequestStatus(id, "denied")
	if err != nil {
		c.String(500, "Failed to reject request")
		return
	}

	c.Redirect(303, "/admin/services/request_info/"+idStr)
}
