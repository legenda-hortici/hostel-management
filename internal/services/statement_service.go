package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type StatementService interface {
	GetAllStatements() ([]models.Statement, error)
	CreateStatementRequest(user_id int, name string, typeStatement string, amount int, request_date string, phone string, hostel int) error
	GetStatementRequestByID(id int) (models.Statement, error)
	UpdateStatementRequestStatus(id int, status string) error
	GetAllUserStatements(email string) ([]models.Statement, error)
}

type statementService struct {
	statementRepo repositories.StatementRepository
}

func NewStatementService(repo repositories.StatementRepository) StatementService {
	return &statementService{
		statementRepo: repo,
	}
}

func (s *statementService) GetAllStatements() ([]models.Statement, error) {
	statements, err := s.statementRepo.GetAllStatements()
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (s *statementService) CreateStatementRequest(user_id int, name string, typeStatement string, amount int, request_date string, phone string, hostel int) error {

	statement := models.Statement{
		Name:     name,
		Type:     typeStatement,
		Amount:   amount,
		Date:     request_date,
		Phone:    phone,
		Status:   "Ожидает",
		Hostel:   hostel,
		Users_id: user_id,
	}

	if statement.Date == "" {
		statement.Date = "Не указана"
	} else if statement.Phone == "" {
		statement.Phone = "Не указан"
	} else if statement.Hostel == 0 {
		statement.Hostel = 0
	}

	if statement.Type == "Бесплатная" {
		statement.Amount = 0
	}

	return s.statementRepo.CreateStatementRequest(statement)
}

func (s *statementService) GetStatementRequestByID(id int) (models.Statement, error) {
	return s.statementRepo.GetStatementRequestByID(id)
}

func (s *statementService) UpdateStatementRequestStatus(id int, status string) error {
	return s.statementRepo.UpdateStatementRequestStatus(id, status)
}

func (s *statementService) GetAllUserStatements(email string) ([]models.Statement, error) {
	return s.statementRepo.GetAllUserStatements(email)
}
