package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/internal/session"
	"hostel-management/storage/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FaqHandler struct {
	faqService services.FaqService
}

func NewFaqHandler(faqService services.FaqService) *FaqHandler {
	return &FaqHandler{
		faqService: faqService,
	}
}

func (h *FaqHandler) SupportHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.SupportHandler"

	role, exists := session.GetUserRole(c)
	if !exists {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v: %v", role, op)
	}

	faq, err := h.faqService.GetAllFaq()
	if err != nil {
		c.String(500, err.Error())
		log.Printf("Failed to get faq: %v: %v", err, op)
	}

	c.HTML(200, "layout.html", gin.H{
		"Page": "support",
		"Role": role,
		"FAQ":  faq,
	})
}

func (h *FaqHandler) AddFaqHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.AddFaqHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v: %v", role, op)
		return
	}

	question := c.PostForm("question")
	answer := c.PostForm("answer")

	faq := models.Faq{
		Question: question,
		Answer:   answer,
	}

	err := h.faqService.CreateFaq(faq)
	if err != nil {
		c.String(500, err.Error())
		log.Printf("Failed to create faq: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/support")
}

func (h *FaqHandler) DeleteFaqHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.DeleteFaqHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		log.Printf("Failed to get ID for faq: %v: %v", err, op)
		return
	}

	err = h.faqService.DeleteFaqItem(id)
	if err != nil {
		c.String(500, "Failed to delete faq")
		log.Printf("Failed to delete faq: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/support")
}

func (h *FaqHandler) UpdateFaqHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.UpdateFaqHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		log.Printf("Failed to get ID for faq: %v: %v", err, op)
		return
	}

	question := c.PostForm("question")
	answer := c.PostForm("answer")

	faq := models.Faq{
		Question: question,
		Answer:   answer,
	}

	err = h.faqService.UpdateFaqItem(id, faq)
	if err != nil {
		c.String(500, "Failed to update faq")
		log.Printf("Failed to update faq: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/support")
}
