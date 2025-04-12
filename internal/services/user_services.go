package services

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
	"time"
)

type UserService interface {
	CreateUser(username, surname, email, password, institute string, roomNumber int) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserRole(email string) (string, error)
	GetTotalCountUsers(searchTerm, filterRole string) (int, error)
	GetPasswordByEmail(email string) (string, error)
	UpdateUser(id int, user *models.UserRequest) error
	UpdateUserByEmail(email, name, surname, password string) error
	DeleteUser(id int) error
	GetResidentsCount() (int, error)
	GetUserIDByEmail(email string) (int, error)
	GetUsernameByID(id int) (string, error)
	GetUserPasswordByEmail(email string) (string, error)
	GetAdminData(role string) (*models.User, error)
	GetHeadmanData(role string) (*models.User, error)
	UpdateAdminData(models.UserRequest) error
	UpdateHeadmanData(models.UserRequest) error
}

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) CreateUser(username, surname, email, password, institute string, roomNumber int) error {
	resident := models.User{
		Username:     username,
		Surname:      surname,
		Email:        email,
		Password:     password,
		Role:         "user",
		Institute:    institute,
		RoomNumber:   roomNumber,
		SettlingDate: time.Now().Format("2006-01-02"),
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
func (s *userServiceImpl) UpdateUser(id int, req *models.UserRequest) error {
	resident := models.User{
		Username:  req.Username,
		Surname:   req.Surname,
		Email:     req.Email,
		Institute: req.Institute,
		Role:      req.Role,
		Password:  req.Password,
	}
	return s.userRepo.Update(id, &resident)
}

// UpdateUserByEmail обновляет данные пользователя
func (s *userServiceImpl) UpdateUserByEmail(email, name, surname, password string) error {
	if email == "" || name == "" || password == "" {
		return errors.New("fields cannot be empty")
	}
	user := &models.User{
		Username: name,
		Surname:  surname,
		Password: password,
	}
	return s.userRepo.UpdateByEmail(email, user)
}

// DeleteUser удаляет пользователя
func (s *userServiceImpl) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
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

func (s *userServiceImpl) GetHeadmanData(role string) (*models.User, error) {
	return s.userRepo.GetHeadman(role)
}

// UpdateAdminData обновляет данные администратора
func (s *userServiceImpl) UpdateAdminData(req models.UserRequest) error {
	admin := models.User{
		Username: req.Username,
		Surname:  req.Surname,
		Password: req.Password,
	}
	return s.userRepo.UpdateAdminData(admin)
}

// UpdateHeadmanData обновляет данные администратора
func (s *userServiceImpl) UpdateHeadmanData(req models.UserRequest) error {
	headman := models.User{
		Username: req.Username,
		Surname:  req.Surname,
		Password: req.Password,
	}
	return s.userRepo.UpdateHeadmanData(headman)
}
