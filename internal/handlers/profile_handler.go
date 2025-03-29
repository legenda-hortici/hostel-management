package handlers

import (
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
	"hostel-management/storage/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userService services.UserService
}

func NewProfileHandler(userService services.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func (h *ProfileHandler) Profile(c *gin.Context) {
	const op = "handlers.ProfileHandler.ProfileHandler"

	role, exists := session.GetUserRole(c)
	if !exists {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v: %v", role, op)
	}

	email, exists := session.GetUserEmail(c)
	if !exists {
		c.String(500, "failed to get user email")
		log.Printf("Failed to get user email: %v: %v", email, op)
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		c.String(500, "failed to get user")
		log.Printf("Failed to get user: %v: %v", err, op)
	}

	c.HTML(200, "layout.html", gin.H{
		"Page": "profile",
		"Role": role,
		"User": user,
	})
}

func (h *ProfileHandler) UpdateProfileHandler(c *gin.Context) {
	const op = "handlers.ProfileHandler.UpdateProfileHandler"

	email, exists := session.GetUserEmail(c)
	if !exists {
		c.String(403, "Access denied")
		log.Printf("Access denied: %v: %v", email, op)
	}

	if c.Request.Method != "POST" {
		log.Printf(" Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	name := c.PostForm("username")
	emailUdp := c.PostForm("email")
	password := c.PostForm("password")

	user := models.User{
		Username: name,
		Email:    emailUdp,
		Password: password,
	}

	err := h.userService.UpdateUserByEmail(email, &user)
	if err != nil {
		log.Printf("Failed to update user: %v: %v", err, op)
		c.String(500, "failed to update user")
		return
	}

	c.Redirect(303, "/profile")
}
