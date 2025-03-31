package main

import (
	"hostel-management/internal/handlers"
	"hostel-management/internal/services"
	"hostel-management/pkg/auth"
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
	inventoryRepo := repositories.NewInventoryRepository()
	faqRepo := repositories.NewFaqRepository()
	hostelRepo := repositories.NewHostelRepository()
	newsRepo := repositories.NewNewsRepository()
	noticeRepo := repositories.NewNoticeRepository()

	authService := auth.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	serviceService := services.NewServiceService(serviceRepo)
	statementService := services.NewStatementService(statementRepo)
	inventoryService := services.NewInventoryService(inventoryRepo)
	faqService := services.NewFaqService(faqRepo)
	hostelService := services.NewHostelService(hostelRepo)
	newsService := services.NewNewsService(newsRepo)
	noticeService := services.NewNoticeService(noticeRepo)

	authHandler := auth.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, roomService)
	adminHandler := handlers.NewAdminHandler(userService, hostelService)
	roomHandler := handlers.NewRoomHandler(roomService)
	serviceHandler := handlers.NewServiceHandler(serviceService, userService, statementService, roomService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	faqhandler := handlers.NewFaqHandler(faqService)
	homeHandler := handlers.NewHomeHandler(newsService, noticeService)
	newsHandler := handlers.NewNewsHandler(newsService)
	noticeHandler := handlers.NewNoticeHandler(noticeService)
	profileHandler := handlers.NewProfileHandler(userService)

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
		protected.GET("/profile", profileHandler.Profile)
		protected.POST("/profile/update_profile", profileHandler.UpdateProfileHandler)
		protected.GET("/services", serviceHandler.ServicesHandler)
		protected.GET("/services/:id", serviceHandler.ServiceInfoHandler)
		protected.POST("/services/send_request/:id", serviceHandler.RequestServiceHandler)
		protected.GET("/services/request_info/:id", serviceHandler.RequestInfoHandler)
		protected.GET("/support", faqhandler.SupportHandler)

		protected.GET("/", homeHandler.HomeHandler)
		protected.POST("/upload_banner", homeHandler.UploadBannerHandler)
		protected.POST("/delete_banner", homeHandler.DeleteBannerHandler)

		protected.GET("/news", newsHandler.News)
		protected.GET("/news/:id", newsHandler.NewsInfoHandler)
		protected.GET("/create_news", newsHandler.CreateNewsPageHandler)
		protected.POST("/create_news/add", newsHandler.CreateNewsHandler)
		protected.POST("/news/:id/delete", newsHandler.DeleteNewsHandler)

		protected.GET("/notices", noticeHandler.Notices)
		protected.GET("/notices/:id", noticeHandler.NoticeInfoHandler)
		protected.GET("/create_notice", noticeHandler.CreateNoticePageHandler)
		protected.POST("/create_notice/add", noticeHandler.CreateNoticeHandler)
		protected.POST("/notices/:id/delete", noticeHandler.DeleteNoticeHandler)
	}

	// Административные маршруты
	admin := r.Group("/admin")
	admin.Use(auth.AdminMiddleware())
	{
		admin.GET("/", adminHandler.AdminCabinetHandler)
		admin.POST("/update_profile", adminHandler.UpdateCabinetHandler)

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
		// TODO : update resident info
		admin.POST("/residents/resident/:id/update_info", userHandler.UpdateResidentDataHandler)

		admin.GET("/services", serviceHandler.ServicesHandler)
		admin.GET("/services/service/:id", serviceHandler.ServiceInfoHandler)
		admin.POST("/services/add_service", serviceHandler.AddServiceHandler)
		admin.POST("/services/service/:id/delete", serviceHandler.DeleteServiceHandler)
		admin.POST("/services/service/:id/edit", serviceHandler.UpdateServiceHandler)
		admin.GET("/services/request_info/:id", serviceHandler.RequestInfoHandler)
		admin.POST("/services/request_info/:id/approve", serviceHandler.AcceptRequestHandler)
		admin.POST("/services/request_info/:id/reject", serviceHandler.RejectRequestHandler)

		admin.GET("/documents", handlers.DocumentsHandler)
		admin.POST("/documents/create_contract", handlers.CreateContractHandler)

		admin.GET("/inventory", inventoryHandler.InventoryHandler)
		admin.POST("/inventory/:id/delete", inventoryHandler.DeleteInventoryItemHandler)
		admin.POST("/inventory/add_item", inventoryHandler.AddInventoryItemHandler)

		admin.GET("/support", faqhandler.SupportHandler)
		admin.POST("/support/add_faq", faqhandler.AddFaqHandler)
		admin.POST("/support/faq/:id/delete", faqhandler.DeleteFaqHandler)
		admin.POST("/support/faq/:id/update", faqhandler.UpdateFaqHandler)

		// admin.GET("/notices", handlers.AdminNoticesHandler)
	}

	return nil
}
