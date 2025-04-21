package routes

import (
	"hostel-management/internal/dependencies"
	"hostel-management/pkg/auth"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(r *gin.RouterGroup, deps *dependencies.Dependencies) {
	r.GET("/login", auth.GuestMiddleware(), deps.AuthService.LoginHandler)
	r.POST("/login", auth.GuestMiddleware(), deps.AuthService.LoginHandler)
	r.GET("/logout", auth.AuthMiddleware(), deps.AuthService.LogoutHandler)
}
