package main

import (
	"context"
	"hostel-management/cmd/server"
	"hostel-management/storage/db"
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

	// Инициализируем Gin
	router := gin.Default()

	// Настраиваем сессии
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("hostel_session", store))

	// Регистрируем маршруты
	if err := RegisterRoutes(router); err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Запускаем сервер
	if err := server.StartServer(ctx, router); err != nil {
		log.Printf("Server error: %v", err)
	}
}
