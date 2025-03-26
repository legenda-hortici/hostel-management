package main

import (
	"hostel-management/internal/auth"
	"hostel-management/internal/handlers"
	"hostel-management/internal/helpers"
	"hostel-management/internal/services"
	"hostel-management/storage/repositories"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) error {
	// Настраиваем шаблоны
	r.LoadHTMLGlob("web/templates/*.html")
	r.Static("/static", "./web/static")

	userRepo := repositories.NewUserRepository()
	roomRepo := repositories.NewRoomRepository()
	serviceRepo := repositories.NewServiceRepository()
	statementRepo := repositories.NewStatementRepository()

	authService := auth.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	serviceService := services.NewServiceService(serviceRepo)
	statementService := services.NewStatementService(statementRepo)

	userHelper := helpers.NewUserHelper()
	roomHelper := helpers.NewRoomHelper()

	authHandler := auth.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, roomService, userHelper)
	roomHandler := handlers.NewRoomHandler(roomService, roomHelper)
	serviceHandler := handlers.NewServiceHandler(serviceService, userService, statementService)

	// Публичные маршруты
	public := r.Group("/")
	{
		public.GET("/login", auth.GuestMiddleware(), authHandler.LoginHandler)
		public.POST("/login", auth.GuestMiddleware(), authHandler.LoginHandler)
		public.GET("/logout", auth.AuthMiddleware(), authHandler.LogoutHandler)
	}

	// Защищенные маршруты
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/profile", handlers.ProfileHandler)
		protected.GET("/", handlers.HomeHandler)
	}

	// Административные маршруты
	admin := r.Group("/admin")
	admin.Use(auth.AdminMiddleware())
	{
		admin.GET("/", userHandler.AdminCabinetHandler)

		admin.GET("/rooms", roomHandler.RoomsHandler)
		admin.POST("/rooms/add_room", roomHandler.AddRoomHandler)
		admin.GET("/rooms/room_info/:id", roomHandler.RoomInfoHandler)
		admin.POST("/rooms/room_info/:id/add_resident_into_room", roomHandler.AddResidentIntoRoomHandler)
		admin.POST("/rooms/room_info/delete_from_room", roomHandler.DeleteResidentFromRoomHandler)
		admin.POST("/rooms/room_info/:id/freeze", roomHandler.FreezeRoomHandler)

		admin.GET("/residents", userHandler.ResidentsHandler)
		admin.GET("/residents/resident/:id", userHandler.ResidentInfoHandler)
		admin.POST("/residents/add_resident", userHandler.AddResidentHandler)
		admin.POST("/residents/:id", userHandler.UpdateResidentDataHandler)
		admin.POST("/residents/resident/:id/delete_resident", userHandler.DeleteResidentHandler)

		admin.GET("/services", serviceHandler.ServicesHandler)

		// admin.GET("/inventory", handlers.AdminInventoryHandler)
		// admin.GET("/news", handlers.AdminNewsHandler)
		// admin.GET("/notices", handlers.AdminNoticesHandler)
		// admin.GET("/documents", handlers.AdminDocumentsHandler)
		// admin.GET("/support", handlers.AdminSupportHandler)
	}

	return nil
}
