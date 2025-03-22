package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"users-test-task-/internal/handlers"
	"users-test-task-/internal/models"
	"users-test-task-/internal/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Mock User Repository

type MockUserRepository struct {
	Users  map[int]*models.User
	LastID int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:  make(map[int]*models.User),
		LastID: 0,
	}
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	if user.Name == "error" {
		return errors.New("intentional error")
	}
	m.LastID++
	user.ID = m.LastID
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUser(id int) (*models.User, error) {
	user, exists := m.Users[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *MockUserRepository) UpdateUser(id int, user *models.User) error {
	_, exists := m.Users[id]
	if !exists {
		return sql.ErrNoRows
	}
	m.Users[id] = user
	return nil
}

func (m *MockUserRepository) DeleteUser(id int) error {
	_, exists := m.Users[id]
	if !exists {
		return sql.ErrNoRows
	}
	delete(m.Users, id)
	return nil
}

// Test User Service

func TestCreateUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := services.NewUserService(repo)

	user := &models.User{Name: "Test", Email: "test@test.com", Phone: "+1234567890"}
	err := service.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
}

func TestGetUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := services.NewUserService(repo)

	user := &models.User{Name: "Test", Email: "test@test.com", Phone: "+1234567890"}
	service.CreateUser(user)

	retrievedUser, err := service.GetUser(1)
	assert.NoError(t, err)
	assert.Equal(t, user, retrievedUser)
}

func TestUpdateUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := services.NewUserService(repo)

	user := &models.User{Name: "Test", Email: "test@test.com", Phone: "+1234567890"}
	service.CreateUser(user)

	updatedUser := &models.User{ID: 1, Name: "Updated", Email: "updated@test.com", Phone: "+0987654321"}
	err := service.UpdateUser(1, updatedUser)
	assert.NoError(t, err)

	retrievedUser, _ := service.GetUser(1)
	assert.Equal(t, "Updated", retrievedUser.Name)
}

func TestDeleteUser(t *testing.T) {
	repo := NewMockUserRepository()
	service := services.NewUserService(repo)

	user := &models.User{Name: "Test", Email: "test@test.com", Phone: "+1234567890"}
	service.CreateUser(user)

	err := service.DeleteUser(1)
	assert.NoError(t, err)

	retrievedUser, err := service.GetUser(1)
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
}

// Test Handlers
func TestCreateUserHandler(t *testing.T) {
	repo := NewMockUserRepository()
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	user := models.User{Name: "Test", Email: "test@test.com", Phone: "+1234567890"}
	jsonBody, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestGetUserHandler(t *testing.T) {
	repo := NewMockUserRepository()
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	user := &models.User{Name: "Test", Email: "test@test.com", Phone: "+1234567890"}
	service.CreateUser(user)

	req, _ := http.NewRequest("GET", "/users/1", nil)

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
