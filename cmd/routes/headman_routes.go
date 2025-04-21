package routes

import (
	"hostel-management/internal/dependencies"

	"github.com/gin-gonic/gin"
)

func RegisterHeadmanRoutes(r *gin.RouterGroup, deps *dependencies.Dependencies) {
	r.GET("/", deps.HeadmanHandler.HeadmanCabinetHandler)
	r.POST("/update_profile", deps.HeadmanHandler.UpdateHeadmanData)
	

	rooms := r.Group("/rooms")
	{
		rooms.GET("/", deps.RoomHandler.RoomsHandler)
		rooms.GET("/room_info/:id", deps.RoomHandler.RoomInfoHandler)
		rooms.GET("/room_info/resident/:id", deps.UserHandler.ResidentInfoHandler)
		rooms.POST("/room_info/:id/add_resident_into_room", deps.RoomHandler.AddResidentIntoRoomHandler)
		rooms.POST("/room_info/delete_from_room", deps.RoomHandler.DeleteResidentFromRoomHandler)
		rooms.POST("/room_info/:id/freeze", deps.RoomHandler.FreezeRoomHandler)
		rooms.POST("/room_info/:id/unfreeze", deps.RoomHandler.UnfreezeRoomHandler)
	}

	residents := r.Group("/residents")
	{
		residents.GET("/", deps.UserHandler.ResidentsHandler)
		residents.GET("/resident/:id", deps.UserHandler.ResidentInfoHandler)
		residents.POST("/add_resident", deps.UserHandler.AddResidentHandler)
		residents.POST("/:id", deps.UserHandler.UpdateResidentDataHandler)
		residents.POST("/resident/:id/delete_resident", deps.UserHandler.DeleteResidentHandler)
		residents.PUT("/resident/:id/edit", deps.UserHandler.UpdateResidentDataHandler)
	}

	services := r.Group("/services")
	{
		services.GET("/", deps.ServiceHandler.ServicesHandler)
		services.GET("/service/:id", deps.ServiceHandler.ServiceInfoHandler)
		services.GET("/request_info/:id", deps.ServiceHandler.RequestInfoHandler)
		services.POST("/request_info/:id/approve", deps.ServiceHandler.AcceptRequestHandler)
		services.POST("/request_info/:id/reject", deps.ServiceHandler.RejectRequestHandler)
	}

	inventory := r.Group("/inventory")
	{
		inventory.GET("/", deps.InventoryHandler.InventoryHandler)
		inventory.POST("/:id/delete", deps.InventoryHandler.DeleteInventoryItemHandler)
		inventory.POST("/add_item", deps.InventoryHandler.AddInventoryItemHandler)
		inventory.POST("/update_item", deps.InventoryHandler.UpdateInventoryItemHandler)
	}
}
