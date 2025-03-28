package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type FaqService interface {
	CreateFaq(faq models.Faq) error
	GetFaqByID(id int) (models.Faq, error)
	GetAllFaq() ([]models.Faq, error)
	DeleteFaqItem(id int) error
	UpdateFaqItem(id int, faq models.Faq) error
}

type faqService struct {
	repo repositories.FaqRepository
}

func NewFaqService(repo repositories.FaqRepository) FaqService {
	return &faqService{
		repo: repo,
	}
}

func (s *faqService) CreateFaq(faq models.Faq) error {
	return s.repo.CreateFaq(faq)
}

func (s *faqService) GetAllFaq() ([]models.Faq, error) {
	return s.repo.GetAllFaq()
}

func (s *faqService) GetFaqByID(id int) (models.Faq, error) {
	return s.repo.GetFaqByID(id)
}

func (s *faqService) DeleteFaqItem(id int) error {
	return s.repo.DeleteFaqItem(id)
}

func (s *faqService) UpdateFaqItem(id int, faq models.Faq) error {
	return s.repo.UpdateFaqItem(id, faq)
}
