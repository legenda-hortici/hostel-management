package services

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type NoticeService interface {
	GetAllNotices() ([]models.Notice, error)
	GetNoticeByID(id int) (models.Notice, error)
	CreateNotice(title, annotation, text, date string) error
	GetLatestNotices() ([]models.Notice, error)
	DeleteNotice(id int) error
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

func (s *noticeService) CreateNotice(title, annotation, text, date string) error {
	if title == "" || annotation == "" || text == "" || date == "" {
		return errors.New("fields cannot be empty")
	}

	notice := models.Notice{
		Title:      title,
		Annotation: annotation,
		Text:       text,
		Date:       date,
	}

	return s.noticeRepo.CreateNotice(notice)
}

func (s *noticeService) GetLatestNotices() ([]models.Notice, error) {
	return s.noticeRepo.GetLatestNotices()
}

func (s *noticeService) DeleteNotice(id int) error {
	return s.noticeRepo.DeleteNotice(id)
}
