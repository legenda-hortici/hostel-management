package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
)

type NoticeRepository interface {
	CreateNotice(notice models.Notice) error
	GetAllNotices() ([]models.Notice, error)
	GetLatestNotices() ([]models.Notice, error)
	GetNoticeByID(id int) (models.Notice, error)
	DeleteNotice(id int) error
}

type noticeRepository struct {
	db *sql.DB
}

func NewNoticeRepository() NoticeRepository {
	return &noticeRepository{
		db: db.DB,
	}
}

func (r *noticeRepository) CreateNotice(notice models.Notice) error {
	query := "INSERT INTO News (title, annotation, text, date, type) VALUES (?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, notice.Title, notice.Annotation, notice.Text, notice.Date, "Срочная")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *noticeRepository) GetAllNotices() ([]models.Notice, error) {
	query := "SELECT * FROM News WHERE type = 'Срочная'"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := []models.Notice{}
	for rows.Next() {
		new := models.Notice{}
		err := rows.Scan(&new.ID, &new.Title, &new.Annotation, &new.Text, &new.Date, &new.NewsType)
		if err != nil {
			return nil, err
		}
		news = append(news, new)
	}

	return news, nil
}

func (r *noticeRepository) GetLatestNotices() ([]models.Notice, error) {
	query := "SELECT * FROM News WHERE type = 'Срочная' ORDER BY date DESC LIMIT 3"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notices := []models.Notice{}
	for rows.Next() {
		new := models.Notice{}
		err := rows.Scan(&new.ID, &new.Title, &new.Annotation, &new.Text, &new.Date, &new.NewsType)
		if err != nil {
			return nil, err
		}
		notices = append(notices, new)
	}

	return notices, nil
}

func (r *noticeRepository) GetNoticeByID(id int) (models.Notice, error) {
	query := "SELECT * FROM News WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	news := models.Notice{}
	err := row.Scan(&news.ID, &news.Title, &news.Annotation, &news.Text, &news.Date, &news.NewsType)
	if err != nil {
		return models.Notice{}, err
	}
	return news, nil
}

func (r *noticeRepository) DeleteNotice(id int) error {
	query := "DELETE FROM News WHERE id = ?"
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
