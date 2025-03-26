package services

import (
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserRole(email string) (string, error)
	GetTotalCountUsers(searchTerm, filterRole string) (int, error)
	GetPasswordByEmail(email string) (string, error)
	UpdateUser(user *models.User) error
	UpdateUserByEmail(email string, user *models.User) error
	DeleteUser(id int) error
	GetAdminInfo(role string) (*models.User, error)
	GetResidentsCount() (int, error)
	GetUserIDByEmail(email string) (int, error)
	GetUsernameByID(id int) (string, error)
	GetUserPasswordByEmail(email string) (string, error)
	GetAdminData(role string) (*models.User, error)
}

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userServiceImpl) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userServiceImpl) GetUserRole(email string) (string, error) {
	return s.userRepo.GetRole(email)
}

func (s *userServiceImpl) GetTotalCountUsers(searchTerm, filterRole string) (int, error) {
	return s.userRepo.GetTotalCount(searchTerm, filterRole)
}

func (s *userServiceImpl) GetPasswordByEmail(email string) (string, error) {
	return s.userRepo.GetPasswordByEmail(email)
}

func (s *userServiceImpl) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *userServiceImpl) UpdateUserByEmail(email string, user *models.User) error {
	return s.userRepo.UpdateByEmail(email, user)
}

func (s *userServiceImpl) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}

func (s *userServiceImpl) GetAdminInfo(role string) (*models.User, error) {
	return s.userRepo.GetAdminInfo(role)
}

func (s *userServiceImpl) GetResidentsCount() (int, error) {
	return s.userRepo.GetResidentsCount()
}

func (s *userServiceImpl) GetUserIDByEmail(email string) (int, error) {
	return s.userRepo.GetUserIDByEmail(email)
}

func (s *userServiceImpl) GetUsernameByID(id int) (string, error) {
	return s.userRepo.GetUsernameByID(id)
}

func (s *userServiceImpl) GetUserPasswordByEmail(email string) (string, error) {
	return s.userRepo.GetUserPasswordByEmail(email)
}

func (s *userServiceImpl) GetAdminData(role string) (*models.User, error) {
	return s.userRepo.GetAdmin(role)
}
