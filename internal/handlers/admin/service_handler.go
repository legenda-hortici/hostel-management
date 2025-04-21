package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/middlewares"
	handlers "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
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
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	services, err := h.serviceService.GetAllServices()
	if err != nil {
		log.Printf("Error getting services: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить услуги")
		return
	}

	var statements []models.Statement
	if role == "admin" {
		statements, err = h.statementService.GetAllStatements()
		if err != nil {
			log.Printf("Error getting statements: %v: %v", err, op)
			middlewares.HandleError(c, 500, "Ошибка: не удалось получить заявки")
			return
		}
	} else if role == "headman" {
		statements, err = h.statementService.GetAllStatementsByHeadman(email)
		if err != nil {
			log.Printf("Error getting statements: %v: %v", err, op)
			middlewares.HandleError(c, 500, "Ошибка: не удалось получить заявки")
			return
		}
	}

	userStatements, err := h.statementService.GetAllUserStatements(email)
	if err != nil {
		log.Printf("Failed to get user statements: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить заявки")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":           "services",
		"Role":           role,
		"Services":       services,
		"Statements":     statements,
		"userStatements": userStatements,
		"Flashes":        flashes,
	})
}

func (h *ServiceHandler) AddServiceHandler(c *gin.Context) {
	const op = "handlers.ServiceHandler.AddServiceHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
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
		middlewares.HandleError(c, 500, "Ошибка: не удалось создать услугу")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Услуга успешно создана!")
	session.Save()

	c.Redirect(303, "/admin/services")
}

func (h *ServiceHandler) ServiceInfoHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.ServiceInfoHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for service: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить ID услуги")
		return
	}

	service, err := h.serviceService.GetServiceByID(id)
	if err != nil {
		log.Printf("Failed to get service: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить услугу")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "service_info",
		"Role":    role,
		"Service": service,
		"Flashes": flashes,
	})
}

func (h *ServiceHandler) UpdateServiceHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.UpdateServiceHandler"

	// Проверяем скрытое поле _method
	if c.Request.Method != "POST" || c.Request.FormValue("_method") != "PUT" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for service: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить ID услуги")
		return
	}

	err = h.serviceService.UpdateServiceByID(id, c)
	if err != nil {
		log.Printf("Failed to update service: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось обновить услугу")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Услуга успешно обновлена!")
	session.Save()

	c.Redirect(303, "/admin/services/service/"+idStr)
}

func (h *ServiceHandler) DeleteServiceHandler(c *gin.Context) {
	const op = "handlers.ServiceHandler.DeleteServiceHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for service: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить ID услуги")
		return
	}

	err = h.serviceService.DeleteService(id)
	if err != nil {
		log.Printf("Failed to delete service: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить услугу")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Услуга удалена!")
	session.Save()

	c.Redirect(303, "/admin/services")
}

func (h *ServiceHandler) RequestServiceHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RequestServiceHandler"

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	_, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid service ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID услуги")
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("Failed to get user: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить пользователя")
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
		middlewares.HandleError(c, 500, "Ошибка: не удалось создать заявку на услугу")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Заявка отправлена!")
	session.Save()

	c.Redirect(303, "/services")
}

func (h *ServiceHandler) RequestInfoHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RequestInfoHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	request, err := h.statementService.GetStatementRequestByID(id)
	if err != nil {
		log.Printf("Failed to get service request: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить заявку на услугу")
		return
	}

	user, err := h.userService.GetUserByID(request.Users_id)
	if err != nil {
		log.Printf("Failed to get user: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить пользователя")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":      "request_info",
		"Role":      role,
		"Statement": request,
		"User":      user,
		"Flashes":   flashes,
	})
}

func (h *ServiceHandler) AcceptRequestHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.AcceptRequestHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	err = h.statementService.UpdateStatementRequestStatus(id, "Одобрена")
	if err != nil {
		log.Printf("Failed to accept request: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось одобрить заявку")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Успешно!")
	session.Save()

	c.Redirect(303, "/admin/services/request_info/"+idStr)
}

func (h *ServiceHandler) RejectRequestHandler(c *gin.Context) {

	const op = "handlers.ServiceHandler.RejectRequestHandler"

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	err = h.statementService.UpdateStatementRequestStatus(id, "Отклонена")
	if err != nil {
		log.Printf("Failed to reject request: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось отклонить заявку")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Успешно!")
	session.Save()

	c.Redirect(303, "/admin/services/request_info/"+idStr)
}
