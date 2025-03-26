package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	UserRoleKey  = "user_role"
)

// CreateSession создает новую сессию пользователя
func CreateSession(c *gin.Context, userID int, email, role string) {
	session := sessions.Default(c)
	session.Set(UserIDKey, userID)
	session.Set(UserEmailKey, email)
	session.Set(UserRoleKey, role)
	session.Save()
}

// DeleteSession удаляет сессию пользователя
func DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// GetUserID получает ID пользователя из сессии
func GetUserID(c *gin.Context) (int, bool) {
	session := sessions.Default(c)
	userID := session.Get(UserIDKey)
	if userID == nil {
		return 0, false
	}
	return userID.(int), true
}

// GetUserEmail получает email пользователя из сессии
func GetUserEmail(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	email := session.Get(UserEmailKey)
	if email == nil {
		return "", false
	}
	return email.(string), true
}

// GetUserRole получает роль пользователя из сессии
func GetUserRole(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	role := session.Get(UserRoleKey)
	if role == nil {
		return "", false
	}
	return role.(string), true
}

// IsAuthenticated проверяет, аутентифицирован ли пользователь
func IsAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	return session.Get(UserIDKey) != nil
}

// IsAdmin проверяет, является ли пользователь администратором
func IsAdmin(c *gin.Context) bool {
	session := sessions.Default(c)
	role := session.Get(UserRoleKey)
	if role == nil {
		return false
	}
	return role.(string) == "admin"
}
