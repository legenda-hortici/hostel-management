package headman

import (
	"hostel-management/internal/services"
	handlers "hostel-management/pkg/validation"

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
	if err != nil {
		c.String(403, err.Error())
		return
	}

	c.HTML(200, "layout.html", gin.H{
		"Page": "headman_cabinet",
		"Role": role,
	})
}
