package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type NoticeService interface {
	GetAllNotices() ([]models.Notice, error)
	GetNoticeByID(id int) (models.Notice, error)
	CreateNotice(notice models.Notice) error
	GetLatestNotices() ([]models.Notice, error)
}

type noticeService struct {
	noticeRepo repositories.NoticeRepository
}

func NewNoticeService(repo repositories.NoticeRepository) NoticeService {
	return &noticeService{
		noticeRepo: repo,
	}
}

func (s *noticeService) GetAllNotices() ([]models.Notice, error) {
	return s.noticeRepo.GetAllNotices()
}

func (s *noticeService) GetNoticeByID(id int) (models.Notice, error) {
	return s.noticeRepo.GetNoticeByID(id)
}

func (s *noticeService) CreateNotice(notice models.Notice) error {
	return s.noticeRepo.CreateNotice(notice)
}

func (s *noticeService) GetLatestNotices() ([]models.Notice, error) {
	return s.noticeRepo.GetLatestNotices()
}
