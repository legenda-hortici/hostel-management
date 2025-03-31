package main

import (
	"context"
	"hostel-management/cmd/routes"
	"hostel-management/cmd/server"
	"hostel-management/internal/config"
	"hostel-management/storage/db"
	"hostel-management/storage/redis"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем обработку сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received shutdown signal: %v", sig)
		cancel() // Отменяем контекст при получении сигнала
	}()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		db.CloseDB()
		log.Println("All resources have been cleaned up")
	}()

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Config loaded successfully")

	// Инициализируем Redis
	redisCache, err := redis.NewRedisCache(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisCache.Close()

	// Инициализируем Gin
	router := gin.Default()

	// Настраиваем сессии
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("hostel_session", store))

	// Регистрируем маршруты
	if err := routes.RegisterRoutes(router, redisCache); err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Запускаем сервер
	if err := server.StartServer(ctx, router); err != nil {
		log.Printf("Server error: %v", err)
	}
}
