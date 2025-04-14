package handlers

import (
	"hostel-management/internal/services"
	handlers "hostel-management/pkg/validation"
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

func NewServiceHandler(
	serviceService services.ServiceService,
	userService services.UserService,
	statementService services.StatementService,
	roomService services.RoomService) *ServiceHandler {
	return &ServiceHandler{
		serviceService:   serviceService,
		userService:      userService,
		statementService: statementService,
		roomService:      roomService,
	}
}

func (h *ServiceHandler) ServicesHandler(c *gin.Context) {
	const op = "handlers.ServiceHandler.ServicesHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	services, err := h.serviceService.GetAllServices()
	if err != nil {
		log.Printf("Error getting services: %v: %v", err, op)
		c.String(500, "Error getting services: "+err.Error())
		return
	}

	var statements []models.Statement
	if role == "admin" {
		statements, err = h.statementService.GetAllStatements()
		if err != nil {
			log.Printf("Error getting statements: %v: %v", err, op)
			c.String(500, "Error getting statements: "+err.Error())
			return
		}
	} else if role == "headman" {
		statements, err = h.statementService.GetAllStatementsByHeadman(email)
		if err != nil {
			log.Printf("Error getting statements: %v: %v", err, op)
			c.String(500, "Error getting statements: "+err.Error())
			return
		}
	}

	userStatements, err := h.statementService.GetAllUserStatements(email)
	if err != nil {
		log.Printf("Failed to get user statements: %v: %v", err, op)
		c.String(500, "Failed to get user statements: "+err.Error())
		return
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
	const op = "handlers.ServiceHandler.AddServiceHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	name := c.PostForm("name")
	typeService := c.PostForm("type")
	description := c.PostForm("description")
	amount, _ := strconv.Atoi(c.PostForm("cost"))
	isDate := c.PostForm("is_date") == "on"
	isHostel := c.PostForm("is_hostel") == "on"
	isPhone := c.PostForm("is_phone") == "on"

	err := h.serviceService.CreateService(name, typeService, description, isDate, isHostel, isPhone, amount)
	if err != nil {
		log.Printf("Failed to create service: %v: %v", err, op)
		c.String(500, "Failed to create service")
		return
	}

	c.Redirect(303, "/admin/services")
}

func (h *ServiceHandler) ServiceInfoHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.ServiceInfoHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for service: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	service, err := h.serviceService.GetServiceByID(id)
	if err != nil {
		log.Printf("Failed to get service: %v: %v", err, op)
		c.String(500, "Failed to get service")
		return
	}

	// log.Println(service)

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":    "service_info",
		"Role":    role,
		"Service": service,
	})
}

func (h *ServiceHandler) UpdateServiceHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.UpdateServiceHandler"

	// Проверяем скрытое поле _method
	if c.Request.Method != "POST" || c.Request.FormValue("_method") != "PUT" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for service: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	err = h.serviceService.UpdateServiceByID(id, c)
	if err != nil {
		log.Printf("Failed to update service: %v: %v", err, op)
		c.String(500, "Failed to update service")
		return
	}

	c.Redirect(303, "/admin/services/service/"+idStr)
}

func (h *ServiceHandler) DeleteServiceHandler(c *gin.Context) {
	const op = "handlers.ServiceHandler.DeleteServiceHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for service: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	err = h.serviceService.DeleteService(id)
	if err != nil {
		log.Printf("Failed to delete service: %v: %v", err, op)
		c.String(500, "Failed to delete service")
		return
	}

	c.Redirect(303, "/admin/services")
}

func (h *ServiceHandler) RequestServiceHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RequestServiceHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	_, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid service ID: %v: %v", err, op)
		c.String(400, "Invalid service ID")
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("Failed to get user: %v: %v", err, op)
		c.String(500, "Failed to get user")
		return
	}
	name := c.PostForm("name")
	typeService := c.PostForm("type")
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	request_date := c.PostForm("request_date")
	hostel, _ := strconv.Atoi(c.PostForm("hostel"))
	phone := c.PostForm("phone")

	err = h.statementService.CreateStatementRequest(user.ID, name, typeService, amount, request_date, phone, hostel)
	if err != nil {
		log.Printf("Failed to create service request: %v: %v", err, op)
		c.String(500, "Failed to create service request")
		return
	}

	c.Redirect(303, "/services")
}

func (h *ServiceHandler) RequestInfoHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RequestInfoHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
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
		log.Printf("Failed to get service request: %v: %v", err, op)
		c.String(500, "Failed to get service request")
		return
	}

	user, err := h.userService.GetUserByID(request.Users_id)
	if err != nil {
		log.Printf("Failed to get user: %v: %v", err, op)
		c.String(500, "Failed to get user")
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":      "request_info",
		"Role":      role,
		"Statement": request,
		"User":      user,
	})
}

func (h *ServiceHandler) AcceptRequestHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.AcceptRequestHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	err = h.statementService.UpdateStatementRequestStatus(id, "Одобрена")
	if err != nil {
		log.Printf("Failed to accept request: %v: %v", err, op)
		c.String(500, "Failed to accept request")
		return
	}

	c.Redirect(303, "/admin/services/request_info/"+idStr)
}

func (h *ServiceHandler) RejectRequestHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RejectRequestHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	err = h.statementService.UpdateStatementRequestStatus(id, "Отклонена")
	if err != nil {
		log.Printf("Failed to reject request: %v: %v", err, op)
		c.String(500, "Failed to reject request")
		return
	}

	c.Redirect(303, "/admin/services/request_info/"+idStr)
}
