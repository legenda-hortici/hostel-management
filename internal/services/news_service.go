package services

import (
	"context"
	"errors"
	"fmt"
	"hostel-management/storage/cache"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
	"log"
	"time"
)

type NewsService interface {
	GetAllNews(ctx context.Context) ([]models.News, error)
	GetNewsByID(id int) (models.News, error)
	CreateNews(ctx context.Context, title, annotation, text, date string) error
	GetLatestNews() ([]models.News, error)
	DeleteNews(ctx context.Context, id int) error
}

type newsService struct {
	repo  repositories.NewsRepository
	cache cache.Cache
}

func NewNewsService(repo repositories.NewsRepository, cache cache.Cache) NewsService {
	return &newsService{
		repo:  repo,
		cache: cache,
	}
}

func (s *newsService) GetAllNews(ctx context.Context) ([]models.News, error) {
	// Пробуем получить из кеша
	var news []models.News
	cacheKey := "news:all"
	err := s.cache.Get(ctx, cacheKey, &news)
	if err == nil {
		for i := range news {
			news[i].Date = news[i].Date[:10]
		}
		return news, nil
	}

	// Если в кеше нет, получаем из базы данных
	news, err = s.repo.GetAllNews()
	if err != nil {
		return nil, fmt.Errorf("failed to get news from repository: %w", err)
	}

	for i := range news {
		news[i].Date = news[i].Date[:10]
	}

	// Сохраняем в кеш на 5 минут
	err = s.cache.Set(ctx, cacheKey, news, 5*time.Minute)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		log.Printf("Failed to cache news: %v\n", err)
	}

	return news, nil
}

func (s *newsService) GetNewsByID(id int) (models.News, error) {
	news, err := s.repo.GetNewsByID(id)
	if err != nil {
		return models.News{}, err
	}

	news.Date = news.Date[:10]

	return news, err
}

func (s *newsService) CreateNews(ctx context.Context, title, annotation, text, date string) error {
	if title == "" || annotation == "" || text == "" || date == "" {
		return errors.New("fields cannot be empty")
	}

	news := models.News{
		Title:      title,
		Annotation: annotation,
		Text:       text,
		Date:       date,
	}

	err := s.repo.CreateNews(news)
	if err != nil {
		return fmt.Errorf("failed to create news: %w", err)
	}

	// Удаляем кеш новостей
	err = s.cache.Delete(ctx, "news:all")
	if err != nil {
		fmt.Printf("Failed to invalidate news cache: %v\n", err)
	}

	return nil
}

func (s *newsService) GetLatestNews() ([]models.News, error) {

	news, err := s.repo.GetLatestNews()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest news: %w", err)
	}

	for i := range news {
		news[i].Date = news[i].Date[:10]
	}
	return news, nil
}

func (s *newsService) DeleteNews(ctx context.Context, id int) error {
	err := s.repo.DeleteNews(id)
	if err != nil {
		return fmt.Errorf("failed to delete news: %w", err)
	}

	err = s.cache.Delete(ctx, "news:all")
	if err != nil {
		fmt.Printf("Failed to invalidate news cache: %v\n", err)
	}
	return s.repo.DeleteNews(id)
}
