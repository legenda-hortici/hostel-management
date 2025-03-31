package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/storage/db"
	"hostel-management/storage/models"
)

type NewsRepository interface {
	CreateNews(news models.News) error
	GetAllNews() ([]models.News, error)
	GetLatestNews() ([]models.News, error)
	GetNewsByID(id int) (models.News, error)
	DeleteNews(id int) error
}

type newsRepository struct {
	db *sql.DB
}

func NewNewsRepository() NewsRepository {
	return &newsRepository{
		db: db.DB,
	}
}

func (r *newsRepository) CreateNews(news models.News) error {
	query := "INSERT INTO News (title, annotation, text, date, type) VALUES (?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, news.Title, news.Annotation, news.Text, news.Date, "regular")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *newsRepository) GetAllNews() ([]models.News, error) {
	query := "SELECT * FROM News WHERE type = 'regular'"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := []models.News{}
	for rows.Next() {
		new := models.News{}
		err := rows.Scan(&new.ID, &new.Title, &new.Annotation, &new.Text, &new.Date, &new.NewsType)
		if err != nil {
			return nil, err
		}
		news = append(news, new)
	}

	return news, nil
}

func (r *newsRepository) GetLatestNews() ([]models.News, error) {
	query := "SELECT * FROM News WHERE type = 'regular' ORDER BY date DESC LIMIT 3"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := []models.News{}
	for rows.Next() {
		new := models.News{}
		err := rows.Scan(&new.ID, &new.Title, &new.Annotation, &new.Text, &new.Date, &new.NewsType)
		if err != nil {
			return nil, err
		}
		news = append(news, new)
	}

	return news, nil
}

func (r *newsRepository) GetNewsByID(id int) (models.News, error) {
	query := "SELECT * FROM News WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	news := models.News{}
	err := row.Scan(&news.ID, &news.Title, &news.Annotation, &news.Text, &news.Date, &news.NewsType)
	if err != nil {
		return models.News{}, err
	}
	return news, nil
}

func (r *newsRepository) DeleteNews(id int) error {
	query := "DELETE FROM News WHERE id = ?"
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
