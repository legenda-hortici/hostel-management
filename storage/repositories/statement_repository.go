package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/internal/config/db"
	"hostel-management/storage/models"
)

type StatementRepository interface {
	GetAllStatements() ([]models.Statement, error)
	CreateStatementRequest(statement models.Statement) error
	GetStatementRequestByID(id int) (models.Statement, error)
	UpdateStatementRequestStatus(id int, status string) error
	GetAllUserStatements(email string) ([]models.Statement, error)
}

type statementRepository struct {
	db *sql.DB
}

func NewStatementRepository() StatementRepository {
	return &statementRepository{
		db: db.DB,
	}
}

func (r *statementRepository) GetAllStatements() ([]models.Statement, error) {
	query := `SELECT s.id, s.name, s.type, s.amount, s.date, s.phone, s.status, s.hostel, u.name 
			FROM Statements s
			JOIN Users u ON s.users_id = u.id;`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statements := []models.Statement{}
	for rows.Next() {
		statement := models.Statement{}
		err := rows.Scan(&statement.ID, &statement.Name, &statement.Type, &statement.Amount, &statement.Date, &statement.Phone, &statement.Status, &statement.Hostel, &statement.Username)
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, nil
}

func (r *statementRepository) CreateStatementRequest(statement models.Statement) error {
	query := "INSERT INTO Statements (name, type, amount, date, phone, status, hostel, users_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, statement.Name, statement.Type, statement.Amount, statement.Date, statement.Phone, statement.Status, statement.Hostel, statement.Users_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *statementRepository) GetStatementRequestByID(id int) (models.Statement, error) {
	query := "SELECT * FROM Statements WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	statement := models.Statement{}
	err := row.Scan(&statement.ID, &statement.Name, &statement.Type, &statement.Amount, &statement.Date, &statement.Phone, &statement.Status, &statement.Users_id, &statement.Hostel)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Statement{}, fmt.Errorf("statement with ID %d not found", id)
		}
		return models.Statement{}, err
	}

	return statement, nil
}

func (r *statementRepository) UpdateStatementRequestStatus(id int, status string) error {
	query := "UPDATE Statements SET status = ? WHERE id = ?"
	_, err := db.DB.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *statementRepository) GetAllUserStatements(email string) ([]models.Statement, error) {
	query := `
			SELECT s.id, s.name, s.type, s.amount, s.date, s.phone, s.status, s.hostel, s.users_id
			FROM Statements s
			JOIN Users u ON s.users_id = u.id
			WHERE u.email = ?;
		`
	rows, err := db.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statements := []models.Statement{}

	for rows.Next() {
		statement := models.Statement{}
		err := rows.Scan(&statement.ID, &statement.Name, &statement.Type, &statement.Amount, &statement.Date, &statement.Phone, &statement.Status, &statement.Hostel, &statement.Users_id)
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, nil
}
