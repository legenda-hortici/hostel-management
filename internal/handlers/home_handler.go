package handlers

import (
	"hostel-management/internal/session"
	"log"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	// Проверяем аутентификацию
	if !session.IsAuthenticated(c) {
		log.Println("User is not authenticated")
		c.Redirect(302, "/login")
		return
	}

	// Получаем роль
	role, exists := session.GetUserRole(c)
	if !exists {
		log.Println("User role not found in session")
		c.Redirect(302, "/login")
		return
	}

	c.HTML(200, "layout.html", gin.H{
		"Page": "home",
		"Role": role,
	})
}
