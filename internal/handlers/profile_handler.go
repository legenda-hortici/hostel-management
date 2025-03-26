package handlers

import (
	"hostel-management/internal/session"

	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context) {
	email, _ := session.GetUserEmail(c)
	role, _ := session.GetUserRole(c)

	c.HTML(200, "layout.html", gin.H{
		"Page":  "profile",
		"Email": email,
		"Role":  role,
	})
}
