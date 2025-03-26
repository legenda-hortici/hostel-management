package auth

import (
	"hostel-management/internal/session"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	authService AuthService
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginHandler обрабатывает страницу входа
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login", gin.H{
			"Page": "login",
		})
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.authService.Login(email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login", gin.H{
			"Page":  "login",
			"Error": err.Error(),
		})
		return
	}
	log.Println(user)

	// Создаем сессию
	session.CreateSession(c, user.ID, user.Email, user.Role)

	// Перенаправляем в зависимости от роли
	switch user.Role {
	case "admin":
		c.Redirect(http.StatusFound, "/admin")
	case "user":
		c.Redirect(http.StatusFound, "/")
	default:
		c.Redirect(http.StatusFound, "/")
	}
}

// LogoutHandler обрабатывает выход пользователя
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	session.DeleteSession(c)
	c.Redirect(http.StatusFound, "/login")
}
