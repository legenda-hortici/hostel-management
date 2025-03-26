package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() error {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
	}

	// Получаем параметры подключения из переменных окружения
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Формируем строку подключения
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Проверяем подключение
	err = DB.Ping()
	if err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	fmt.Println("Успешное подключение к базе данных")

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
