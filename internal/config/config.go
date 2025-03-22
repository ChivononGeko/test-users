package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	AppPort    string
}

func LoadConfig() (*Config, error) {
	slog.Info("Загрузка конфигурации из .env файла...")

	err := godotenv.Load()
	if err != nil {
		slog.Error("Ошибка при загрузке .env файла", "error", err)
		return nil, err
	}

	config := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		AppPort:    os.Getenv("APP_PORT"),
	}

	slog.Info("Конфигурация успешно загружена.")
	return config, nil
}

// GetDBURL формирует URL для подключения к базе данных
func (cfg *Config) GetDBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
}
