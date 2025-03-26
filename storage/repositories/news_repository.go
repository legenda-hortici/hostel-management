package repositories

import (
	"fmt"
	"hostel-management/storage/db"
	"hostel-management/storage/models"
)

func CreateNews(news models.News) error {
	query := "INSERT INTO News (title, annotation, text, date, type) VALUES (?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, news.Title, news.Annotation, news.Text, news.Date, "regular")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func CreateNotice(notice models.News) error {
	query := "INSERT INTO News (title, annotation, text, date, type) VALUES (?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, notice.Title, notice.Annotation, notice.Text, notice.Date, "breaking")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetAllNews() ([]models.News, error) {
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

func GetLatestNews() ([]models.News, error) {
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

func GetAllNotices() ([]models.News, error) {
	query := "SELECT * FROM News WHERE type = 'breaking'"
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

func GetLatestNotices() ([]models.News, error) {
	query := "SELECT * FROM News WHERE type = 'breaking' ORDER BY date DESC LIMIT 3"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notices := []models.News{}
	for rows.Next() {
		new := models.News{}
		err := rows.Scan(&new.ID, &new.Title, &new.Annotation, &new.Text, &new.Date, &new.NewsType)
		if err != nil {
			return nil, err
		}
		notices = append(notices, new)
	}

	return notices, nil
}

func GetNewsByID(id int) (models.News, error) {
	query := "SELECT * FROM News WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	news := models.News{}
	err := row.Scan(&news.ID, &news.Title, &news.Annotation, &news.Text, &news.Date, &news.NewsType)
	if err != nil {
		return models.News{}, err
	}
	return news, nil
}

func GetNoticeByID(id int) (models.News, error) {
	query := "SELECT * FROM News WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	news := models.News{}
	err := row.Scan(&news.ID, &news.Title, &news.Annotation, &news.Text, &news.Date, &news.NewsType)
	if err != nil {
		return models.News{}, err
	}
	return news, nil
}
