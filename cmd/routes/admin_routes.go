package routes

import (
	"hostel-management/internal/dependencies"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, deps *dependencies.Dependencies) {
	r.GET("/", deps.AdminHandler.AdminCabinetHandler)
	r.POST("/update_profile", deps.AdminHandler.UpdateCabinetHandler)
	r.POST("/create_contract", deps.AdminHandler.CreateContractHandler)
	// TODO: Сделать выгрузку общей статистики по общежитиям

	hostel := r.Group("/hostel")
	{
		hostel.GET("/:id", deps.AdminHandler.HostelInfoHandler)
		hostel.POST(":id/assign_headman", deps.AdminHandler.AssignCommandantHandler)
		hostel.POST(":id/remove_headman", deps.AdminHandler.RemoveCommandantHandler)
	}

	rooms := r.Group("/rooms")
	{
		rooms.GET("/", deps.RoomHandler.RoomsHandler)
		rooms.POST("/add_room", deps.RoomHandler.AddRoomHandler)
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
		services.POST("/add_service", deps.ServiceHandler.AddServiceHandler)
		services.POST("/service/:id/delete", deps.ServiceHandler.DeleteServiceHandler)
		services.POST("/service/:id/edit", deps.ServiceHandler.UpdateServiceHandler)
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

	support := r.Group("/support")
	{
		support.GET("/", deps.FaqHandler.SupportHandler)
		support.POST("/add_faq", deps.FaqHandler.AddFaqHandler)
		support.POST("/faq/:id/delete", deps.FaqHandler.DeleteFaqHandler)
		support.POST("/faq/:id/update", deps.FaqHandler.UpdateFaqHandler)
	}
}
