package service

import (
	"errors"
	"strings"

	"users/internal/models"
	"users/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(name, email string) (*models.User, error)
	UpdateUser(id int, name, email string) (*models.User, error)
	DeleteUser(id int) error
}

type userService struct {
	repo repository.UserRepository
}

// Inject Repository ผ่าน Constructor ตามหลัก DIP
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) CreateUser(name, email string) (*models.User, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
		return nil, errors.New("name and email are required")
	}
	user, err := s.repo.Create(name, email)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, errors.New("email already exists")
	}
	return user, err
}

func (s *userService) UpdateUser(id int, name, email string) (*models.User, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
		return nil, errors.New("name and email are required")
	}
	user, err := s.repo.Update(id, name, email)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, errors.New("email already exists")
	}
	return user, err
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
