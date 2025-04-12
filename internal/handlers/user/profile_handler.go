package handlers

import (
	"hostel-management/internal/services"
	handlers "hostel-management/pkg/validation"
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

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
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

	email, err := handlers.ValidateUserByEmail(c, op)
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
	surname := c.PostForm("surname")
	password := c.PostForm("password")

	err = h.userService.UpdateUserByEmail(email, name, surname, password)
	if err != nil {
		log.Printf("Failed to update user: %v: %v", err, op)
		c.String(500, "failed to update user")
		return
	}

	c.Redirect(303, "/profile")
}
