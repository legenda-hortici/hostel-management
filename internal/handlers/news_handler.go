package handlers

import (
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"hostel-management/storage/models"
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

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" && role != "user" {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v", op)
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

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v", op)
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

	news := models.News{
		Title:      title,
		Annotation: annotation,
		Text:       text,
		Date:       date,
	}

	err := h.newsService.CreateNews(news)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/")
}

func (h *NewsHandler) NewsInfoHandler(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.String(405, "NewsInfoHandler: Method not allowed")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(400, "NewsInfoHandler: Некорректный ID")
		return
	}

	news, err := h.newsService.GetNewsByID(id)
	if err != nil {
		c.String(500, "NewsInfoHandler: Failed to get news")
		return
	}

	news.NewsType = helpers.TranslateNewsType(news.NewsType)
	news.Date = news.Date[:10]

	c.HTML(200, "layout.html", map[string]interface{}{
		"Page": "news_info",
		"News": news,
	})
}
