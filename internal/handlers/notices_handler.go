package handlers

import (
	"errors"
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
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

func ValidateUserByRole(c *gin.Context, op string) (string, error) {
	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" && role != "user" {
		log.Printf("Access denied: %v", op)
		return "", errors.New("access denied")
	}
	return role, nil
}

func (h *NoticesHandler) Notices(c *gin.Context) {

	const op = "handlers.NoticesHandler.Notices"

	role, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
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

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":        "notices",
		"Role":        role,
		"Notices":     notices,
		"LatestNotes": latestNotes,
	})
}

func (h *NoticesHandler) CreateNoticePageHandler(c *gin.Context) {
	data := map[string]interface{}{
		"Page": "create_notices",
	}
	c.HTML(200, "layout.html", data)
}

func (h *NoticesHandler) CreateNoticeHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.CreateNoticeHandler"

	_, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, "Access denied")
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

	err = h.noticeService.CreateNotice(title, annotation, text, date)
	if err != nil {
		log.Printf("Failed to create notice: %v: %v", err, op)
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/")
}

func (h *NoticesHandler) NoticeInfoHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.NoticeInfoHandler"

	_, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
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

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":    "notices_info",
		"Notices": notice,
	})
}

func (h *NoticesHandler) DeleteNoticeHandler(c *gin.Context) {

	const op = "handlers.NoticesHandler.DeleteNoticeHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		log.Printf("Failed to get ID for notice: %v: %v", err, op)
		return
	}

	err = h.noticeService.DeleteNotice(id)
	if err != nil {
		c.String(500, "Failed to delete notice")
		log.Printf("Failed to delete notice: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/notices")
}
