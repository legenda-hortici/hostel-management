package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type StatementService interface {
	GetAllStatements() ([]models.Statement, error)
	CreateStatementRequest(statement models.Statement) error
	GetStatementRequestByID(id int) (models.Statement, error)
	UpdateStatementRequestStatus(id int, status string) error
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
	return s.statementRepo.GetAllStatements()
}

func (s *statementService) CreateStatementRequest(statement models.Statement) error {
	return s.statementRepo.CreateStatementRequest(statement)
}

func (s *statementService) GetStatementRequestByID(id int) (models.Statement, error) {
	return s.statementRepo.GetStatementRequestByID(id)
}

func (s *statementService) UpdateStatementRequestStatus(id int, status string) error {
	return s.statementRepo.UpdateStatementRequestStatus(id, status)
}
