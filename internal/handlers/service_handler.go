package handlers

import (
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/internal/session"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceService   services.ServiceService
	userService      services.UserService
	statementService services.StatementService
	serviceHelper    helpers.ServiceHelper
}

func NewServiceHandler(serviceService services.ServiceService, userService services.UserService, statementService services.StatementService) *ServiceHandler {
	return &ServiceHandler{
		serviceService:   serviceService,
		userService:      userService,
		statementService: statementService,
	}
}

func (h *ServiceHandler) ServicesHandler(c *gin.Context) {
	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied")
		return
	}

	services, err := h.serviceService.GetAllServices()
	if err != nil {
		c.String(500, "Ошибка получения сервисов: "+err.Error())
		return
	}

	statements, err := h.statementService.GetAllStatements()
	if err != nil {
		c.String(500, "Ошибка получения заявок: "+err.Error())
		return
	}
	log.Println(statements)

	for statement := range statements {
		statements[statement].Status = helpers.TranslateStatus(statements[statement].Status)

	}

	data := map[string]interface{}{
		"Page":       "services",
		"Role":       role,
		"Services":   services,
		"Statements": statements,
	}
	c.HTML(200, "layout.html", data)
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

	c.Redirect(303, "/services")
}

func (h *ServiceHandler) ServiceInfoHandler(c *gin.Context) {
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
		"Service": service,
	}
	c.HTML(200, "layout.html", data)
}

func (h *ServiceHandler) UpdateServiceHandler(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		return
	}

	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.String(400, "Invalid request body")
		return
	}

	err = h.serviceService.UpdateServiceByID(id, service)
	if err != nil {
		c.String(500, "Failed to update service")
		return
	}

	c.Status(200)
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

	c.Redirect(303, "/services")
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

	statement := models.Statement{
		Name:     c.PostForm("name"),
		Type:     c.PostForm("type"),
		Amount:   0,
		Date:     "",
		Phone:    "",
		Status:   "awaits",
		Hostel:   0,
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
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
		c.String(500, "Failed to get user")
		return
	}

	data := map[string]interface{}{
		"Page":    "request_info",
		"Request": request,
		"User":    user,
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

	c.Redirect(303, "/services")
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

	c.Redirect(303, "/services")
}
