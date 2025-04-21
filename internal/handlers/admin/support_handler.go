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

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	faq, err := h.faqService.GetAllFaq()
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить FAQ")
		log.Printf("Failed to get faq: %v: %v", err, op)
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "support",
		"Role":    role,
		"FAQ":     faq,
		"Flashes": flashes,
	})
}

func (h *FaqHandler) AddFaqHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.AddFaqHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Access denied: %v", err)
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	question := c.PostForm("question")
	answer := c.PostForm("answer")

	err = h.faqService.CreateFaq(question, answer)
	if err != nil {
		log.Printf("Failed to create faq: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось создать FAQ")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Успешно!")
	session.Save()

	c.Redirect(303, "/admin/support")
}

func (h *FaqHandler) DeleteFaqHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.DeleteFaqHandler"

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for faq: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить ID FAQ")
		return
	}

	err = h.faqService.DeleteFaqItem(id)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить FAQ")
		log.Printf("Failed to delete faq: %v: %v", err, op)
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Успешно!")
	session.Save()

	c.Redirect(303, "/admin/support")
}

func (h *FaqHandler) UpdateFaqHandler(c *gin.Context) {

	const op = "handlers.SupportHandler.UpdateFaqHandler"

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить ID FAQ")
		log.Printf("Failed to get ID for faq: %v: %v", err, op)
		return
	}

	question := c.PostForm("question")
	answer := c.PostForm("answer")

	err = h.faqService.UpdateFaqItem(id, question, answer)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось обновить FAQ")
		log.Printf("Failed to update faq: %v: %v", err, op)
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Успешно обновлено!")
	session.Save()

	c.Redirect(303, "/admin/support")
}
