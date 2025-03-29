package auth

import (
	"hostel-management/pkg/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет аутентификацию пользователя
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !session.IsAuthenticated(c) {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminMiddleware проверяет, является ли пользователь администратором
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !session.IsAdmin(c) {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GuestMiddleware проверяет, что пользователь не аутентифицирован
func GuestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session.IsAuthenticated(c) {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
