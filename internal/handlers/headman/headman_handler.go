package headman

import (
	"hostel-management/internal/config/db"
	"hostel-management/internal/services"
	handlers "hostel-management/pkg/validation"
	"hostel-management/storage/models"
	"log"

	"github.com/gin-gonic/gin"
)

type HeadmanHandler struct {
	userService   services.UserService
	hostelService services.HostelService
}

func NewHeadmanHandler(userService services.UserService, hostelService services.HostelService) *HeadmanHandler {
	return &HeadmanHandler{
		userService:   userService,
		hostelService: hostelService,
	}
}

func (h *HeadmanHandler) HeadmanCabinetHandler(c *gin.Context) {

	const op = "handlers.headman.HeadmanHandler.HeadmanCabinetHandler"

	role, err := handlers.ValidateUserByRole(c, op)
	if err != nil && role != "headman" {
		log.Printf("access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	email, err := handlers.ValidateUserByEmail(c, op)
	if err != nil {
		log.Printf("access denied: %v", err)
		c.String(403, err.Error())
		return
	}

	headmanData, err := h.userService.GetHeadmanData(role)
	if err != nil {
		log.Printf("failed to get headman data: %v", err)
		c.String(500, err.Error()+": "+op)
		return
	}

	hostelData, err := h.hostelService.GetHostelInfoByHeadman(db.DB, email)
	if err != nil {
		log.Printf("failed to get hostel info: %v: %v", err, op)
		c.String(500, err.Error()+": "+op)
		return
	}

	log.Println(hostelData)

	c.HTML(200, "layout.html", gin.H{
		"Page": "headman_cabinet",
		"Role": role,
		"Headman": map[string]interface{}{
			"Username": headmanData.Username,
			"Surname":  headmanData.Surname,
			"Password": headmanData.Password,
			"Email":    headmanData.Email,
		},
		"Hostel": hostelData,
	})
}

func (h *HeadmanHandler) UpdateHeadmanData(c *gin.Context) {

	const op = "handlers.headman.UpdateHeadmanData"

	_, err := handlers.ValidateUserByRole(c, op)
	if err != nil {
		log.Printf("Failed to update headman data: %v: %v", err, op)
		c.String(403, err.Error())
		return
	}

	username := c.PostForm("username")
	surname := c.PostForm("surname")
	password := c.PostForm("password")

	var req models.UserRequest

	req.Username = username
	req.Surname = surname
	req.Password = password

	err = h.userService.UpdateHeadmanData(req)
	if err != nil {
		log.Printf("Failed to update headman data: %v: %v", err, op)
		c.String(500, err.Error()+": "+op)
		return
	}

	c.Redirect(303, "/headman")
}
