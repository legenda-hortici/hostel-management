package cache

import (
	"context"
	"time"
)

// Cache интерфейс для работы с кешем
type Cache interface {
	// Set сохраняет данные в кеш с указанным временем жизни
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get получает данные из кеша
	Get(ctx context.Context, key string, dest interface{}) error

	// Delete удаляет данные из кеша
	Delete(ctx context.Context, key string) error

	// Clear очищает все данные из кеша
	Clear(ctx context.Context) error

	// Close закрывает соединение с кешем
	Close() error
}
