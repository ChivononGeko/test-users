package db

import (
	"database/sql"
	"log/slog"
	"users-test-task-/internal/config"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(cfg *config.Config) (*sql.DB, error) {
	slog.Info("Установка соединения с базой данных...")

	dbURL := cfg.GetDBURL()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Ошибка при подключении к базе данных", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("Ошибка при проверке соединения с базой данных", "error", err)
		return nil, err
	}

	slog.Info("Соединение с базой данных успешно установлено.")
	return db, nil
}
