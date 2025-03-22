package repositories

import (
	"database/sql"
	"log/slog"
	"time"
	"users-test-task-/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	slog.Info("Создание пользователя", "user", user)

	query := `INSERT INTO users (name, email, phone, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, user.Name, user.Email, user.Phone, time.Now()).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetUser(id int) (*models.User, error) {
	slog.Info("Получение пользователя", "userID", id)

	query := `SELECT id, name, email, phone, created_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(id int, user *models.User) error {
	slog.Info("Обновление пользователя", "userID", user.ID)

	query := `UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, id)
	return err
}

func (r *UserRepository) DeleteUser(id int) error {
	slog.Info("Удаление пользователя", "userID", id)

	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
