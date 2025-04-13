package routes

import (
	adm "hostel-management/internal/handlers/admin"
	head "hostel-management/internal/handlers/headman"
	user "hostel-management/internal/handlers/user"
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
	userHandler := adm.NewUserHandler(userService, roomService)
	adminHandler := adm.NewAdminHandler(userService, hostelService)
	headmanHandler := head.NewHeadmanHandler(userService, hostelService)
	roomHandler := adm.NewRoomHandler(roomService)
	serviceHandler := adm.NewServiceHandler(serviceService, userService, statementService, roomService)
	inventoryHandler := adm.NewInventoryHandler(inventoryService)
	faqhandler := adm.NewFaqHandler(faqService)
	homeHandler := adm.NewHomeHandler(newsService, noticeService)
	newsHandler := adm.NewNewsHandler(newsService)
	noticeHandler := adm.NewNoticeHandler(noticeService)
	profileHandler := user.NewProfileHandler(userService)

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

	// TODO: Сделать обработчики для коменданта
	headman := r.Group("/headman")
	headman.Use(auth.HeadmanMiddleware())
	{
		headman.GET("/", headmanHandler.HeadmanCabinetHandler)
		headman.POST("/update_profile", headmanHandler.UpdateHeadmanData)

		rooms := headman.Group("/rooms")
		{
			rooms.GET("/", roomHandler.RoomsHandler)
			rooms.GET("/room_info/:id", roomHandler.RoomInfoHandler)
			rooms.GET("/room_info/resident/:id", userHandler.ResidentInfoHandler)
			rooms.POST("/room_info/:id/add_resident_into_room", roomHandler.AddResidentIntoRoomHandler)
			rooms.POST("/room_info/delete_from_room", roomHandler.DeleteResidentFromRoomHandler)
			rooms.POST("/room_info/:id/freeze", roomHandler.FreezeRoomHandler)
		}

		residents := headman.Group("/residents")
		{
			residents.GET("/", userHandler.ResidentsHandler)
			residents.GET("/resident/:id", userHandler.ResidentInfoHandler)
			residents.POST("/add_resident", userHandler.AddResidentHandler)
			residents.POST("/:id", userHandler.UpdateResidentDataHandler)
			residents.POST("/resident/:id/delete_resident", userHandler.DeleteResidentHandler)
			residents.PUT("/resident/:id/edit", userHandler.UpdateResidentDataHandler)
		}

		services := headman.Group("/services")
		{
			services.GET("/", serviceHandler.ServicesHandler)
			services.GET("/service/:id", serviceHandler.ServiceInfoHandler)
			services.GET("/request_info/:id", serviceHandler.RequestInfoHandler)
			services.POST("/request_info/:id/approve", serviceHandler.AcceptRequestHandler)
			services.POST("/request_info/:id/reject", serviceHandler.RejectRequestHandler)
		}

		inventory := headman.Group("/inventory")
		{
			inventory.GET("/", inventoryHandler.InventoryHandler)
			inventory.POST("/:id/delete", inventoryHandler.DeleteInventoryItemHandler)
			inventory.POST("/add_item", inventoryHandler.AddInventoryItemHandler)
		}
	}

	// TODO: Сделать обработку ошибок с помощью JS
	// TODO: Сделать фотографию в профиле у пользователей

	// Административные маршруты
	admin := r.Group("/admin")
	admin.Use(auth.AdminMiddleware())
	{
		admin.GET("/", adminHandler.AdminCabinetHandler)
		admin.POST("/update_profile", adminHandler.UpdateCabinetHandler)
		admin.POST("/create_contract", adm.CreateContractHandler)
		// TODO: Сделать выгрузку общей статистики по общежитиям

		hostel := admin.Group("/hostel")
		{
			hostel.GET("/:id", adminHandler.HostelInfoHandler)
			hostel.POST(":id/assign_headman", adminHandler.AssignCommandantHandler)
		}

		rooms := admin.Group("/rooms")
		{
			rooms.GET("/", roomHandler.RoomsHandler)
			rooms.POST("/add_room", roomHandler.AddRoomHandler)
			rooms.GET("/room_info/:id", roomHandler.RoomInfoHandler)
			rooms.GET("/room_info/resident/:id", userHandler.ResidentInfoHandler)
			rooms.POST("/room_info/:id/add_resident_into_room", roomHandler.AddResidentIntoRoomHandler)
			rooms.POST("/room_info/delete_from_room", roomHandler.DeleteResidentFromRoomHandler)
			rooms.POST("/room_info/:id/freeze", roomHandler.FreezeRoomHandler)
			rooms.POST("/room_info/:id/unfreeze", roomHandler.UnfreezeRoomHandler)
		}

		residents := admin.Group("/residents")
		{
			residents.GET("/", userHandler.ResidentsHandler)
			residents.GET("/resident/:id", userHandler.ResidentInfoHandler)
			residents.POST("/add_resident", userHandler.AddResidentHandler)
			residents.POST("/:id", userHandler.UpdateResidentDataHandler)
			residents.POST("/resident/:id/delete_resident", userHandler.DeleteResidentHandler)
			residents.PUT("/resident/:id/edit", userHandler.UpdateResidentDataHandler)
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

		inventory := admin.Group("/inventory")
		{
			inventory.GET("/", inventoryHandler.InventoryHandler)
			inventory.POST("/:id/delete", inventoryHandler.DeleteInventoryItemHandler)
			inventory.POST("/add_item", inventoryHandler.AddInventoryItemHandler)
			inventory.POST("/update_item", inventoryHandler.UpdateInventoryItemHandler)
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
