package services

import (
	"users-test-task-/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUser(id int) (*models.User, error)
	UpdateUser(id int, user *models.User) error
	DeleteUser(id int) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) UpdateUser(id int, user *models.User) error {
	return s.repo.UpdateUser(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
