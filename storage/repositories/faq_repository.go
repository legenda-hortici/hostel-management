package repositories

import (
	"hostel-management/storage/db"
	"hostel-management/storage/models"
)

func CreateFaq(faq models.Faq) error {
	query := "INSERT INTO Faq (question, answer) VALUES (?, ?)"
	_, err := db.DB.Exec(query, faq.Question, faq.Answer)
	if err != nil {
		return err
	}
	return nil
}

func GetAllFaq() ([]models.Faq, error) {
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

func GetFaqByID(id int) (models.Faq, error) {
	query := "SELECT * FROM Faq WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	faq := models.Faq{}
	err := row.Scan(&faq.Question, &faq.Answer)
	if err != nil {
		return models.Faq{}, err
	}
	return faq, nil
}
