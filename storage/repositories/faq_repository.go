package repositories

import (
	"database/sql"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
)

type FaqRepository interface {
	CreateFaq(faq models.Faq) error
	GetAllFaq() ([]models.Faq, error)
	GetFaqByID(id int) (models.Faq, error)
	DeleteFaqItem(id int) error
	UpdateFaqItem(id int, faq models.Faq) error
}

type faqRepository struct {
	db *sql.DB
}

func NewFaqRepository() FaqRepository {
	return &faqRepository{
		db: db.DB,
	}
}

func (r *faqRepository) CreateFaq(faq models.Faq) error {
	query := "INSERT INTO Faq (question, answer) VALUES (?, ?)"
	_, err := db.DB.Exec(query, faq.Question, faq.Answer)
	if err != nil {
		return err
	}
	return nil
}

func (r *faqRepository) GetAllFaq() ([]models.Faq, error) {
	query := "SELECT * FROM Faq"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	faqs := []models.Faq{}
	for rows.Next() {
		faq := models.Faq{}
		err := rows.Scan(&faq.ID, &faq.Question, &faq.Answer)
		if err != nil {
			return nil, err
		}
		faqs = append(faqs, faq)
	}

	return faqs, nil
}

func (r *faqRepository) GetFaqByID(id int) (models.Faq, error) {
	query := "SELECT * FROM Faq WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	faq := models.Faq{}
	err := row.Scan(&faq.Question, &faq.Answer)
	if err != nil {
		return models.Faq{}, err
	}
	return faq, nil
}

func (r *faqRepository) DeleteFaqItem(id int) error {
	query := "DELETE FROM Faq WHERE id = ?"
	_, err := db.DB.Exec(query, id)
	return err
}

func (r *faqRepository) UpdateFaqItem(id int, faq models.Faq) error {
	query := "UPDATE Faq SET question = ?, answer = ? WHERE id = ?"
	_, err := db.DB.Exec(query, faq.Question, faq.Answer, id)
	return err
}
