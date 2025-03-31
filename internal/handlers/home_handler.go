package handlers

import (
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"log"
	"os"
	"path/filepath"

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

	news, err := h.newsService.GetAllNews(c)
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

func (h *HomeHandler) UploadBannerHandler(c *gin.Context) {

	const op = "handlers.HomeHandler.UploadBannerHandler"
	const bannerDir = "web/static/banners/"

	// Проверяем аутентификацию
	if !session.IsAuthenticated(c) {
		log.Printf("User is not authenticated: %v", op)
		c.Redirect(302, "/login")
		return
	}

	// Получаем роль
	_, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	file, err := c.FormFile("banner")
	if err != nil {
		log.Printf("Error uploading file: %v: %v", err, op)
		c.String(400, "Ошибка загрузки")
		return
	}

	filePath := filepath.Join(bannerDir, file.Filename)
	log.Printf("File path: %v", filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v: %v", err, op)
		c.String(500, "Ошибка сохранения файла")
		return
	}

	c.Redirect(303, "/")
}

func (h *HomeHandler) DeleteBannerHandler(c *gin.Context) {

	const op = "handlers.HomeHandler.DeleteBannerHandler"

	// Проверяем аутентификацию
	if !session.IsAuthenticated(c) {
		log.Printf("User is not authenticated: %v", op)
		c.Redirect(302, "/login")
		return
	}

	// Получаем роль
	_, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Метод не разрешен")
		return
	}

	bannerName := c.PostForm("banner")
	if bannerName == "" {
		log.Printf("Banner name not specified: %v", op)
		c.String(400, "name not specified")
		return
	}

	bannerPath := filepath.Join("web", bannerName)
	if err := os.Remove(bannerPath); err != nil {
		log.Printf("Error deleting banner: %v: %v", err, op)
		c.String(500, "error deleting banner: "+err.Error())
		return
	}

	log.Printf("Banner deleted: %s", bannerName)
	c.Redirect(303, "/")
}
