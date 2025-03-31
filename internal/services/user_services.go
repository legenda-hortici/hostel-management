package services

import (
	"database/sql"
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type UserService interface {
	CreateUser(username, email, password, institute string, roomNumber int) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserRole(email string) (string, error)
	GetTotalCountUsers(searchTerm, filterRole string) (int, error)
	GetPasswordByEmail(email string) (string, error)
	UpdateUser(user *models.User) error
	UpdateUserByEmail(email, name, emailUdp, password string) error
	DeleteUser(id int) error
	GetAdminInfo(role string) (*models.User, error)
	GetResidentsCount() (int, error)
	GetUserIDByEmail(email string) (int, error)
	GetUsernameByID(id int) (string, error)
	GetUserPasswordByEmail(email string) (string, error)
	GetAdminData(role string) (*models.User, error)
	UpdateAdminData(username, password string) error
}

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) CreateUser(username, email, password, institute string, roomNumber int) error {
	resident := models.User{
		Username:   username,
		Email:      email,
		Password:   password,
		Role:       "user",
		Institute:  sql.NullString{String: institute, Valid: institute != ""},
		RoomNumber: roomNumber,
	}
	return s.userRepo.Create(&resident)
}

// GetAllUsers возвращает список всех пользователей
func (s *userServiceImpl) GetAllUsers() ([]models.User, error) {
	residents, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range residents {
		residents[i].Number = i + 1
	}
	return residents, nil
}

// GetUserByID возвращает пользователя по его ID
func (s *userServiceImpl) GetUserByID(id int) (*models.User, error) {

	return s.userRepo.GetByID(id)
}

// GetUserByEmail возвращает пользователя по его email
func (s *userServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	return s.userRepo.GetByEmail(email)
}

// GetUserRole возвращает роль пользователя
func (s *userServiceImpl) GetUserRole(email string) (string, error) {
	return s.userRepo.GetRole(email)
}

// GetTotalCountUsers возвращает общее количество пользователей
func (s *userServiceImpl) GetTotalCountUsers(searchTerm, filterRole string) (int, error) {
	return s.userRepo.GetTotalCount(searchTerm, filterRole)
}

// GetPasswordByEmail возвращает пароль пользователя
func (s *userServiceImpl) GetPasswordByEmail(email string) (string, error) {
	return s.userRepo.GetPasswordByEmail(email)
}

// UpdateUser обновляет данные пользователя
func (s *userServiceImpl) UpdateUser(req *models.User) error {
	resident := models.User{
		Username: req.Username,
		Email:    req.Email,
		Institute: sql.NullString{
			String: req.Institute.String,
			Valid:  req.Institute.Valid,
		},
		Role:     req.Role,
		Password: req.Password,
	}
	return s.userRepo.Update(&resident)
}

// UpdateUserByEmail обновляет данные пользователя
func (s *userServiceImpl) UpdateUserByEmail(email, name, emailUdp, password string) error {
	if email == emailUdp {
		return s.userRepo.UpdateByEmail(email, &models.User{
			Username: name,
			Password: password,
		})
	}
	if email == "" || emailUdp == "" || name == "" || password == "" {
		return errors.New("fields cannot be empty")
	}
	user := &models.User{
		Username: name,
		Email:    emailUdp,
		Password: password,
	}
	return s.userRepo.UpdateByEmail(email, user)
}

// DeleteUser удаляет пользователя
func (s *userServiceImpl) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}

// GetAdminInfo возвращает информацию о администраторе
func (s *userServiceImpl) GetAdminInfo(role string) (*models.User, error) {
	return s.userRepo.GetAdminInfo(role)
}

// GetResidentsCount возвращает количество жильцов
func (s *userServiceImpl) GetResidentsCount() (int, error) {
	return s.userRepo.GetResidentsCount()
}

// GetUserIDByEmail возвращает ID пользователя по его email
func (s *userServiceImpl) GetUserIDByEmail(email string) (int, error) {
	return s.userRepo.GetUserIDByEmail(email)
}

// GetUsernameByID возвращает имя пользователя по его ID
func (s *userServiceImpl) GetUsernameByID(id int) (string, error) {
	return s.userRepo.GetUsernameByID(id)
}

// GetUserPasswordByEmail возвращает пароль пользователя по его email
func (s *userServiceImpl) GetUserPasswordByEmail(email string) (string, error) {
	return s.userRepo.GetUserPasswordByEmail(email)
}

// GetAdminData возвращает данные администратора
func (s *userServiceImpl) GetAdminData(role string) (*models.User, error) {
	return s.userRepo.GetAdmin(role)
}

// UpdateAdminData обновляет данные администратора
func (s *userServiceImpl) UpdateAdminData(username, password string) error {
	return s.userRepo.UpdateAdminData(username, password)
}
