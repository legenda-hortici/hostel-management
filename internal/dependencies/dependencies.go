package dependencies

import (
	adm "hostel-management/internal/handlers/admin"
	head "hostel-management/internal/handlers/headman"
	user "hostel-management/internal/handlers/user"
	"hostel-management/internal/services"
	"hostel-management/pkg/auth"
	"hostel-management/storage/redis"
	"hostel-management/storage/repositories"
)

type Dependencies struct {
	UserHandler      adm.UserHandler
	HeadmanHandler   head.HeadmanHandler
	ProfileHandler   user.ProfileHandler
	AdminHandler     adm.AdminHandler
	HomeHandler      adm.HomeHandler
	AuthService      auth.AuthHandler
	RoomHandler      adm.RoomHandler
	ServiceHandler   adm.ServiceHandler
	InventoryHandler adm.InventoryHandler
	FaqHandler       adm.FaqHandler
	NewsHandler      adm.NewsHandler
	NoticeHandler    adm.NoticesHandler
}

func BuildDependencies(redisCache *redis.RedisCache) *Dependencies {
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
	headman := head.NewHeadmanHandler(userService, hostelService)
	userHandler := adm.NewUserHandler(userService, roomService)
	profileHandler := user.NewProfileHandler(userService)
	adminHandler := adm.NewAdminHandler(userService, hostelService)
	homeHandler := adm.NewHomeHandler(newsService, noticeService)
	roomHandler := adm.NewRoomHandler(roomService)
	serviceHandler := adm.NewServiceHandler(serviceService, userService, statementService, roomService)
	inventoryHandler := adm.NewInventoryHandler(inventoryService)
	faqHandler := adm.NewFaqHandler(faqService)
	newsHandler := adm.NewNewsHandler(newsService)
	noticeHandler := adm.NewNoticeHandler(noticeService)

	return &Dependencies{
		UserHandler:      *userHandler,
		HeadmanHandler:   *headman,
		ProfileHandler:   *profileHandler,
		AdminHandler:     *adminHandler,
		AuthService:      *authHandler,
		RoomHandler:      *roomHandler,
		ServiceHandler:   *serviceHandler,
		InventoryHandler: *inventoryHandler,
		FaqHandler:       *faqHandler,
		NewsHandler:      *newsHandler,
		NoticeHandler:    *noticeHandler,
		HomeHandler:      *homeHandler,
	}
}
