package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/helpers"
	"hostel-management/pkg/middlewares"
	"hostel-management/pkg/session"
	handlers "hostel-management/pkg/validation"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
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

	if !session.IsAuthenticated(c) {
		log.Printf("User is not authenticated: %v", op)
		c.Redirect(302, "/login")
		return
	}

	role, exists := session.GetUserRole(c)
	if !exists {
		log.Printf("User role not found in session: %v", op)
		c.Redirect(302, "/login")
		return
	}

	news, err := h.newsService.GetAllNews(c)
	if err != nil {
		log.Printf("Error getting news: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить новости")
		return
	}

	notices, err := h.noticesService.GetAllNotices()
	if err != nil {
		log.Printf("Error getting notices: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить уведомления")
		return
	}

	date := time.Now().Format("02.01.2006")

	banners := helpers.GetBanners()

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "home",
		"Date":    date,
		"Role":    role,
		"Banners": banners,
		"News":    news,
		"Notices": notices,
		"Flashes": flashes,
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
	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	file, err := c.FormFile("banner")
	if err != nil {
		log.Printf("Error uploading file: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось загрузить файл")
		return
	}

	filePath := filepath.Join(bannerDir, file.Filename)
	log.Printf("File path: %v", filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось сохранить файл")
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
	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	if c.Request.Method != "POST" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 303, "Ошибка: метод не разрешен")
		return
	}

	bannerName := c.PostForm("banner")
	if bannerName == "" {
		log.Printf("Banner name not specified: %v", op)
		middlewares.HandleError(c, 303, "Ошибка: имя баннера не указано")
		return
	}

	bannerPath := filepath.Join("web", bannerName)
	if err := os.Remove(bannerPath); err != nil {
		log.Printf("Error deleting banner: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить баннер")
		return
	}

	c.Redirect(303, "/")
}
