package services

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type NewsService interface {
	GetAllNews() ([]models.News, error)
	GetNewsByID(id int) (models.News, error)
	CreateNews(title, annotation, text, date string) error
	GetLatestNews() ([]models.News, error)
	DeleteNews(id int) error
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
	news, err := s.repo.GetNewsByID(id)
	if err != nil {
		return models.News{}, err
	}

	switch news.NewsType {
	case "regular":
		news.NewsType = "Регулярная"
	case "breaking":
		news.NewsType = "Срочная"
	default:
		news.NewsType = "Неизвестно"
	}

	news.Date = news.Date[:10]

	return news, err
}

func (s *newsService) CreateNews(title, annotation, text, date string) error {
	if title == "" || annotation == "" || text == "" || date == "" {
		return errors.New("fields cannot be empty")
	}

	news := models.News{
		Title:      title,
		Annotation: annotation,
		Text:       text,
		Date:       date,
	}

	return s.repo.CreateNews(news)
}

func (s *newsService) GetLatestNews() ([]models.News, error) {
	return s.repo.GetLatestNews()
}

func (s *newsService) DeleteNews(id int) error {
	return s.repo.DeleteNews(id)
}
