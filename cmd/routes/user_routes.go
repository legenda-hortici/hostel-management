package routes

import (
	"hostel-management/internal/dependencies"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, deps *dependencies.Dependencies) {
	profile := r.Group("/profile")
	{
		profile.GET("/", deps.ProfileHandler.Profile)
		profile.POST("/update_profile", deps.ProfileHandler.UpdateProfileHandler)
	}

	services := r.Group("/services")
	{
		services.GET("/", deps.ServiceHandler.ServicesHandler)
		services.GET("/:id", deps.ServiceHandler.ServiceInfoHandler)
		services.POST("/send_request/:id", deps.ServiceHandler.RequestServiceHandler)
		services.GET("/request_info/:id", deps.ServiceHandler.RequestInfoHandler)
	}
}
