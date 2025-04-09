package handlers

import (
	"errors"
	"hostel-management/pkg/session"
	"log"

	"github.com/gin-gonic/gin"
)

func ValidateUserByRole(c *gin.Context, op string) (string, error) {
	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" && role != "user" && role != "headman" {
		log.Printf("Access denied: %v", op)
		return "", errors.New("access denied")
	}
	return role, nil
}

func ValidateUserByEmail(c *gin.Context, op string) (string, error) {
	email, exists := session.GetUserEmail(c)
	if !exists {
		log.Printf("Access denied: %v: %v", email, op)
		return "", errors.New("access denied")
	}
	return email, nil
}
