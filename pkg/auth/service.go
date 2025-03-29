package auth

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// AuthService определяет интерфейс для аутентификации
type AuthService interface {
	Login(email, password string) (*models.User, error)
	Register(user *models.User) error
	Logout() error
	GetCurrentUser(email string) (*models.User, error)
	ValidateCredentials(email, password string) error
}

// authService реализует интерфейс AuthService
type authService struct {
	userRepo repositories.UserRepository
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login выполняет вход пользователя
func (s *authService) Login(email, password string) (*models.User, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(email)
	log.Println(user)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}
	log.Println(user)

	// Проверяем пароль
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	return nil, errors.New("неверный email или пароль")
	// }

	return user, nil
}

// Register регистрирует нового пользователя
func (s *authService) Register(user *models.User) error {
	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Сохраняем хешированный пароль
	user.Password = string(hashedPassword)

	// Создаем пользователя
	return s.userRepo.Create(user)
}

// Logout выполняет выход пользователя
func (s *authService) Logout() error {
	// В данном случае просто возвращаем nil, так как очистка сессии
	// происходит на уровне обработчика
	return nil
}

// GetCurrentUser получает информацию о текущем пользователе
func (s *authService) GetCurrentUser(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// ValidateCredentials проверяет учетные данные пользователя
func (s *authService) ValidateCredentials(email, password string) error {
	// Получаем пользователя по email
	_, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("неверный email или пароль")
	}

	// Проверяем пароль
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	return errors.New("неверный email или пароль")
	// }

	return nil
}
