package repositories

import (
	"database/sql"
	"fmt"
	"hostel-management/storage/db"
	"hostel-management/storage/models"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetRole(email string) (string, error)
	GetAll() ([]models.User, error)
	GetTotalCount(searchTerm, filterRole string) (int, error)
	GetPasswordByEmail(email string) (string, error)
	Update(user *models.User) error
	UpdateByEmail(email string, user *models.User) error
	Delete(id int) error
	GetAdminInfo(role string) (*models.User, error)
	GetResidentsCount() (int, error)
	GetUserIDByEmail(email string) (int, error)
	GetUsernameByID(id int) (string, error)
	GetUserPasswordByEmail(email string) (string, error)
	GetAdmin(role string) (*models.User, error)
}

// userRepository реализует интерфейс UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{
		db: db.DB,
	}
}

// Create создает нового пользователя
func (r *userRepository) Create(user *models.User) error {
	query := "INSERT INTO Users (name, email, password, institute, role, Rooms_id) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Institute, user.Role, user.Room_id)
	return err
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := "SELECT id, name, email, password, institute, role, Rooms_id FROM users WHERE id = ?"

	row := r.db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Institute, &user.Role, &user.Room_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден: id %d", user.ID)
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, name, email, role, institute, Rooms_id FROM Users WHERE email = ?"
	row := r.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Institute, &user.Room_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %v", err)
	}

	return user, nil
}

// GetRole получает роль пользователя по email
func (r *userRepository) GetRole(email string) (string, error) {
	var role string
	err := r.db.QueryRow("SELECT role FROM Users WHERE email = ?", email).Scan(&role)
	return role, err
}

// GetAll получает всех пользователей
func (r *userRepository) GetAll() ([]models.User, error) {
	query := "SELECT * FROM Users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Institute, &user.Role, &user.Room_id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetTotalCount получает общее количество пользователей для пагинации
func (r *userRepository) GetTotalCount(searchTerm, filterRole string) (int, error) {
	query := "SELECT COUNT(*) FROM Users WHERE username LIKE ? AND role LIKE ?"
	var count int
	err := r.db.QueryRow(query, "%"+searchTerm+"%", "%"+filterRole+"%").Scan(&count)
	return count, err
}

// GetPasswordByEmail получает пароль пользователя по email
func (r *userRepository) GetPasswordByEmail(email string) (string, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM Users WHERE email = ?", email).Scan(&password)
	return password, err
}

// Update обновляет данные пользователя
func (r *userRepository) Update(user *models.User) error {
	_, err := r.db.Exec(`
        UPDATE Users 
        SET name=?, email=?, institute=?, role=?, password=?
        WHERE id=?`,
		user.Username, user.Email, user.Institute.String, user.Role, user.Password, user.ID,
	)
	return err
}

// UpdateByEmail обновляет данные пользователя по email
func (r *userRepository) UpdateByEmail(email string, user *models.User) error {
	_, err := r.db.Exec(`
		UPDATE Users 
		SET name=?, email=?, password=?
		WHERE email=?`,
		user.Username, user.Email, user.Password, email,
	)
	return err
}

// Delete удаляет пользователя по ID
func (r *userRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM Users WHERE id = ?", id)
	return err
}

// GetAdminInfo получает информацию об администраторе
func (r *userRepository) GetAdminInfo(role string) (*models.User, error) {
	query := "SELECT id, name, email, role FROM Users WHERE role = ?"
	row := r.db.QueryRow(query, role)

	admin := &models.User{}
	err := row.Scan(&admin.ID, &admin.Username, &admin.Email, &admin.Role)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// GetResidentsCount получает количество всех жильцов
func (r *userRepository) GetResidentsCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
	return count, err
}

// GetUserIDByEmail получает ID пользователя по email
func (r *userRepository) GetUserIDByEmail(email string) (int, error) {
	var userID int
	err := r.db.QueryRow("SELECT id FROM Users WHERE email = ?", email).Scan(&userID)
	return userID, err
}

// GetUsernameByID получает имя пользователя по ID
func (r *userRepository) GetUsernameByID(id int) (string, error) {
	var username string
	err := r.db.QueryRow("SELECT name FROM Users WHERE id = ?", id).Scan(&username)
	return username, err
}

// GetUserPasswordByEmail получает пароль пользователя по email
func (r *userRepository) GetUserPasswordByEmail(email string) (string, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM Users WHERE email = ?", email).Scan(&password)
	return password, err
}

func (r *userRepository) GetAdmin(role string) (*models.User, error) {
	query := "SELECT id, name, email, role FROM Users WHERE role = ?"
	row := r.db.QueryRow(query, role)

	admin := &models.User{}
	err := row.Scan(&admin.ID, &admin.Username, &admin.Email, &admin.Role)
	return admin, err
}
