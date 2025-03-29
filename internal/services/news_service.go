package services

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type NewsService interface {
	GetAllNews() ([]models.News, error)
	GetNewsByID(id int) (models.News, error)
	CreateNews(news models.News) error
	GetLatestNews() ([]models.News, error)
}

type newsService struct {
	repo repositories.NewsRepository
}

func NewNewsService(repo repositories.NewsRepository) NewsService {
	return &newsService{
		repo: repo,
	}
}

func (s *newsService) GetAllNews() ([]models.News, error) {
	return s.repo.GetAllNews()
}

func (s *newsService) GetNewsByID(id int) (models.News, error) {
	return s.repo.GetNewsByID(id)
}

func (s *newsService) CreateNews(news models.News) error {
	if news.Title == "" || news.Annotation == "" || news.Text == "" || news.Date == "" {
		return errors.New("fields cannot be empty")
	}
	return s.repo.CreateNews(news)
}

func (s *newsService) GetLatestNews() ([]models.News, error) {
	return s.repo.GetLatestNews()
}
