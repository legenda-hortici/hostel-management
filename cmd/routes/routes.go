package routes

import (
	"hostel-management/internal/handlers"
	"hostel-management/internal/services"
	"hostel-management/pkg/auth"
	"hostel-management/storage/redis"
	"hostel-management/storage/repositories"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, redisCache *redis.RedisCache) error {
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
	newsService := services.NewNewsService(newsRepo, redisCache)
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
		profile := protected.Group("/profile")
		{
			profile.GET("/", profileHandler.Profile)
			profile.POST("/update_profile", profileHandler.UpdateProfileHandler)
		}

		services := protected.Group("/services")
		{
			services.GET("/", serviceHandler.ServicesHandler)
			services.GET("/:id", serviceHandler.ServiceInfoHandler)
			services.POST("/send_request/:id", serviceHandler.RequestServiceHandler)
			services.GET("/request_info/:id", serviceHandler.RequestInfoHandler)
		}

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

		rooms := admin.Group("/rooms")
		{
			rooms.GET("/", roomHandler.RoomsHandler)
			rooms.POST("/add_room", roomHandler.AddRoomHandler)
			rooms.GET("/room_info/:id", roomHandler.RoomInfoHandler)
			rooms.POST("/room_info/:id/add_resident_into_room", roomHandler.AddResidentIntoRoomHandler)
			rooms.POST("/room_info/delete_from_room", roomHandler.DeleteResidentFromRoomHandler)
			rooms.POST("/room_info/:id/freeze", roomHandler.FreezeRoomHandler)
		}

		residents := admin.Group("/residents")
		{
			residents.GET("/", userHandler.ResidentsHandler)
			residents.GET("/resident/:id", userHandler.ResidentInfoHandler)
			residents.POST("/add_resident", userHandler.AddResidentHandler)
			residents.POST("/:id", userHandler.UpdateResidentDataHandler)
			residents.POST("/resident/:id/delete_resident", userHandler.DeleteResidentHandler)
			residents.POST("/resident/:id/update_info", userHandler.UpdateResidentDataHandler)
		}

		services := admin.Group("/services")
		{
			services.GET("/", serviceHandler.ServicesHandler)
			services.GET("/service/:id", serviceHandler.ServiceInfoHandler)
			services.POST("/add_service", serviceHandler.AddServiceHandler)
			services.POST("/service/:id/delete", serviceHandler.DeleteServiceHandler)
			services.POST("/service/:id/edit", serviceHandler.UpdateServiceHandler)
			services.GET("/request_info/:id", serviceHandler.RequestInfoHandler)
			services.POST("/request_info/:id/approve", serviceHandler.AcceptRequestHandler)
			services.POST("/request_info/:id/reject", serviceHandler.RejectRequestHandler)
		}

		documents := admin.Group("/documents")
		{
			documents.GET("/", handlers.DocumentsHandler)
			documents.POST("/create_contract", handlers.CreateContractHandler)
		}

		inventory := admin.Group("/inventory")
		{
			inventory.GET("/", inventoryHandler.InventoryHandler)
			inventory.POST("/:id/delete", inventoryHandler.DeleteInventoryItemHandler)
			inventory.POST("/add_item", inventoryHandler.AddInventoryItemHandler)
		}

		support := admin.Group("/support")
		{
			support.GET("/", faqhandler.SupportHandler)
			support.POST("/add_faq", faqhandler.AddFaqHandler)
			support.POST("/faq/:id/delete", faqhandler.DeleteFaqHandler)
			support.POST("/faq/:id/update", faqhandler.UpdateFaqHandler)
		}
	}

	return nil
}
