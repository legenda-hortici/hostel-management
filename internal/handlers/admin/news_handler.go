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

type NewsHandler struct {
	newsService services.NewsService
}

func NewNewsHandler(newsService services.NewsService) *NewsHandler {
	return &NewsHandler{
		newsService: newsService,
	}
}

func (h *NewsHandler) News(c *gin.Context) {

	const op = "handlers.NewsHandler.NewsHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	news, err := h.newsService.GetAllNews(c)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка получения новостей")
		return
	}

	latestNews, err := h.newsService.GetLatestNews()
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка получения последних новостей")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":       "news",
		"Role":       role,
		"News":       news,
		"LatestNews": latestNews,
		"Flashes":    flashes,
	})
}

func (h *NewsHandler) CreateNewsPageHandler(c *gin.Context) {

	const op = "handlers.CreateNewsPageHandler.CreateNewsPageHandler"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		middlewares.HandleError(c, 403, "Ошибка: доступ запрещен")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", gin.H{
		"Page":    "create_news",
		"Flashes": flashes,
	})
}

func (h *NewsHandler) CreateNewsHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	title := c.PostForm("title")
	annotation := c.PostForm("annotation")
	text := c.PostForm("text")
	date := c.PostForm("date")

	err := h.newsService.CreateNews(c, title, annotation, text, date)
	if err != nil {
		middlewares.HandleError(c, 400, "Ошибка: не удалось создать новость")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Новость успешно создана!")
	session.Save()

	c.Redirect(303, "/")
}

func (h *NewsHandler) NewsInfoHandler(c *gin.Context) {
	const op = "handlers.NewsInfoHandler.NewsInfoHandler"

	if c.Request.Method != "GET" {
		log.Printf("Method not allowed: %v", op)
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for news: %v: %v", err, op)
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		return
	}

	news, err := h.newsService.GetNewsByID(id)
	if err != nil {
		log.Printf("Failed to get news: %v: %v", err, op)
		middlewares.HandleError(c, 500, "Ошибка: не удалось получить новость")
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":    "news_info",
		"News":    news,
		"Flashes": flashes,
	})
}

func (h *NewsHandler) DeleteNewsHandler(c *gin.Context) {

	const op = "handlers.DeleteNewsHandler.DeleteNewsHandler"

	if c.Request.Method != "POST" {
		middlewares.HandleError(c, 405, "Ошибка: метод не разрешен")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middlewares.HandleError(c, 400, "Ошибка: неверный ID")
		log.Printf("Failed to get ID for news: %v: %v", err, op)
		return
	}

	err = h.newsService.DeleteNews(c, id)
	if err != nil {
		middlewares.HandleError(c, 500, "Ошибка: не удалось удалить новость")
		log.Printf("Failed to delete news: %v: %v", err, op)
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Новость успешно удалена!")
	session.Save()

	c.Redirect(303, "/news")
}
