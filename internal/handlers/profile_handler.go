package handlers

import (
	"errors"
	"hostel-management/internal/services"
	"hostel-management/pkg/session"
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

func ValidateUserByEmail(c *gin.Context, op string) (string, error) {
	email, exists := session.GetUserEmail(c)
	if !exists {
		log.Printf("Access denied: %v: %v", email, op)
		return "", errors.New("access denied")
	}
	return email, nil
}

func (h *ProfileHandler) Profile(c *gin.Context) {
	const op = "handlers.ProfileHandler.ProfileHandler"

	role, err := ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	email, err := ValidateUserByEmail(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
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

	email, err := ValidateUserByEmail(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	if c.Request.Method != "POST" {
		log.Printf(" Method not allowed: %v", op)
		c.String(405, "Method not allowed")
		return
	}

	name := c.PostForm("username")
	emailUdp := c.PostForm("email")
	password := c.PostForm("password")

	err = h.userService.UpdateUserByEmail(email, name, emailUdp, password)
	if err != nil {
		log.Printf("Failed to update user: %v: %v", err, op)
		c.String(500, "failed to update user")
		return
	}

	c.Redirect(303, "/profile")
}
