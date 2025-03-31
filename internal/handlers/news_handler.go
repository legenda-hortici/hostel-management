package handlers

import (
	"hostel-management/internal/services"
	"log"
	"strconv"

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

	role, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	news, err := h.newsService.GetAllNews()
	if err != nil {
		c.String(500, "Ошибка получения новостей: "+err.Error())
		return
	}

	latestNews, err := h.newsService.GetLatestNews()
	if err != nil {
		c.String(500, "Ошибка получения последних новостей: "+err.Error())
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page":       "news",
		"Role":       role,
		"News":       news,
		"LatestNews": latestNews,
	})
}

func (h *NewsHandler) CreateNewsPageHandler(c *gin.Context) {

	const op = "handlers.CreateNewsPageHandler.CreateNewsPageHandler"

	_, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page": "create_news",
	})
}

func (h *NewsHandler) CreateNewsHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "CreateNewsHandler: Method not allowed")
		return
	}

	title := c.PostForm("title")
	annotation := c.PostForm("annotation")
	text := c.PostForm("text")
	date := c.PostForm("date")

	err := h.newsService.CreateNews(title, annotation, text, date)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/")
}

func (h *NewsHandler) NewsInfoHandler(c *gin.Context) {
	const op = "handlers.NewsInfoHandler.NewsInfoHandler"

	if c.Request.Method != "GET" {
		log.Printf("Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Failed to get ID for news: %v: %v", err, op)
		c.String(400, "Invalid ID")
		return
	}

	news, err := h.newsService.GetNewsByID(id)
	if err != nil {
		log.Printf("Failed to get news: %v: %v", err, op)
		c.String(500, "Failed to get news")
		return
	}

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page": "news_info",
		"News": news,
	})
}

func (h *NewsHandler) DeleteNewsHandler(c *gin.Context) {

	const op = "handlers.DeleteNewsHandler.DeleteNewsHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		log.Printf("Method not allowed: %v", op)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "Invalid ID")
		log.Printf("Failed to get ID for news: %v: %v", err, op)
		return
	}

	err = h.newsService.DeleteNews(id)
	if err != nil {
		c.String(500, "Failed to delete news")
		log.Printf("Failed to delete news: %v: %v", err, op)
		return
	}

	c.Redirect(303, "/news")
}
