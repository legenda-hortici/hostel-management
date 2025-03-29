package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"hostel-management/storage/models"
	"log"
	"strconv"

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

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" && role != "user" {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v", op)
		return
	}

	notices, err := h.noticeService.GetAllNotices()
	if err != nil {
		log.Printf("Error getting notices: %v: %v", err, op)
		c.String(500, "Ошибка получения объявлений: "+err.Error())
		return
	}

	latestNotes, err := h.noticeService.GetLatestNotices()
	if err != nil {
		log.Printf("Error getting latest notices: %v: %v", err, op)
		c.String(500, "Ошибка получения последних объявлений: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"Page":        "notices",
		"Role":        role,
		"Notices":     notices,
		"LatestNotes": latestNotes,
	}
	c.HTML(200, "layout.html", data)
}

func (h *NoticesHandler) CreateNoticePageHandler(c *gin.Context) {
	data := map[string]interface{}{
		"Page": "create_notices",
	}
	c.HTML(200, "layout.html", data)
}

func (h *NoticesHandler) CreateNoticeHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.CreateNoticeHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "access denied")
		log.Printf("Access denied: %v", op)
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "method not allowed")
		return
	}

	title := c.PostForm("title")
	annotation := c.PostForm("annotation")
	text := c.PostForm("text")
	date := c.PostForm("date")

	notice := models.Notice{
		Title:      title,
		Annotation: annotation,
		Text:       text,
		Date:       date,
	}

	err := h.noticeService.CreateNotice(notice)
	if err != nil {
		log.Printf("Failed to create notice: %v: %v", err, op)
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/")
}

func (h *NoticesHandler) NoticeInfoHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.NoticeInfoHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" && role != "user" {
		c.String(403, "access denied")
		log.Printf("Access denied: %v", op)
		return
	}

	if c.Request.Method != "GET" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for notice: %v: %v", err, op)
		c.String(400, "invalid id")
		return
	}

	notice, err := h.noticeService.GetNoticeByID(id)
	if err != nil {
		log.Printf("Failed to get notice: %v: %v", err, op)
		c.String(500, "failed to get notice")
		return
	}

	data := map[string]interface{}{
		"Page":    "notices_info",
		"Notices": notice,
	}
	c.HTML(200, "layout.html", data)
}
