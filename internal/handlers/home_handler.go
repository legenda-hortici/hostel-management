package handlers

import (
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"log"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	newsService    services.NewsService
	noticesService services.NoticeService
}

func NewHomeHandler(newsService services.NewsService, noticesService services.NoticeService) HomeHandler {
	return HomeHandler{
		newsService:    newsService,
		noticesService: noticesService,
	}
}

func (h *HomeHandler) HomeHandler(c *gin.Context) {

	const op = "handlers.HomeHandler.HomeHandler"

	// Проверяем аутентификацию
	if !session.IsAuthenticated(c) {
		log.Printf("User is not authenticated: %v", op)
		c.Redirect(302, "/login")
		return
	}

	// Получаем роль
	role, exists := session.GetUserRole(c)
	if !exists {
		log.Printf("User role not found in session: %v", op)
		c.Redirect(302, "/login")
		return
	}

	news, err := h.newsService.GetAllNews()
	if err != nil {
		log.Printf("Error getting news: %v: %v", err, op)
		c.String(500, "Error getting news: "+err.Error())
		return
	}

	notices, err := h.noticesService.GetAllNotices()
	if err != nil {
		log.Printf("Error getting notices: %v: %v", err, op)
		c.String(500, "Error getting notices: "+err.Error())
		return
	}

	banners := helpers.GetBanners()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "home",
		"Role":    role,
		"Banners": banners,
		"News":    news,
		"Notices": notices,
	})
}
