package handlers

import (
	"hostel-management/internal/services"
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

	role, err := ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
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

	_, err := ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	question := c.PostForm("question")
	answer := c.PostForm("answer")

	err = h.faqService.CreateFaq(question, answer)
	if err != nil {
		log.Printf("Failed to create faq: %v: %v", err, op)
		c.String(500, err.Error())
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
		log.Printf("Failed to get ID for faq: %v: %v", err, op)
		c.String(400, "Invalid ID")
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

	err = h.faqService.UpdateFaqItem(id, question, answer)
	if err != nil {
		c.String(500, "Failed to update faq")
		log.Printf("Failed to update faq: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/admin/support")
}
