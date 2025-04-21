package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/middlewares"
	handlers "hostel-management/pkg/validation"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type NoticesHandler struct {
	noticeService services.NoticeService
}

func NewNoticeHandler(noticeService services.NoticeService) NoticesHandler {
	return NoticesHandler{
		noticeService: noticeService,
	}
}

func (h *NoticesHandler) Notices(c *gin.Context) {

	const op = "handlers.NoticesHandler.Notices"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	notices, err := h.noticeService.GetAllNotices()
	if err != nil {
		log.Printf("Error getting notices: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка получения объявлений")
		return
	}

	latestNotes, err := h.noticeService.GetLatestNotices()
	if err != nil {
		log.Printf("Error getting latest notices: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка получения последних объявлений")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":        "notices",
		"Role":        role,
		"Notices":     notices,
		"LatestNotes": latestNotes,
		"Flashes":     flashes,
	})
}

func (h *NoticesHandler) CreateNoticePageHandler(c *gin.Context) {

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "create_notices",
		"Flashes": flashes,
	})
}

func (h *NoticesHandler) CreateNoticeHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.CreateNoticeHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, "Access denied")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	title := c.PostForm("title")
	annotation := c.PostForm("annotation")
	text := c.PostForm("text")
	date := c.PostForm("date")

	err = h.noticeService.CreateNotice(title, annotation, text, date)
	if err != nil {
		log.Printf("Failed to create notice: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: не удалось создать объявление")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Объявление успешно создано!")
	session.Save()

	c.Redirect(303, "/")
}

func (h *NoticesHandler) NoticeInfoHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.NoticeInfoHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "GET" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for notice: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	notice, err := h.noticeService.GetNoticeByID(id)
	if err != nil {
		log.Printf("Failed to get notice: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить объявление")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":    "notices_info",
		"Notices": notice,
		"Flashes": flashes,
	})
}

func (h *NoticesHandler) DeleteNoticeHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.DeleteNoticeHandler"

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		log.Printf("Failed to get ID for notice: %v: %v", err, op)
		return
	}

	err = h.noticeService.DeleteNotice(id)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить объявление")
		log.Printf("Failed to delete notice: %v: %v", err, op)
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Объявление успешно удалено!")
	session.Save()

	c.Redirect(303, "/notices")
}
