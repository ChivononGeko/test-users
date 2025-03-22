package main

import (
	"net/http"
	"os"
	"users-test-task-/internal/config"
	"users-test-task-/internal/db"
	"users-test-task-/internal/handlers"
	"users-test-task-/internal/repositories"
	"users-test-task-/internal/services"

	"log/slog"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Ошибка загрузки конфигурации", "error", err)
		os.Exit(1)
	}

	database, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		slog.Error("Ошибка подключения к базе данных", "error", err)
		os.Exit(1)
	}

	if err := db.RunMigrations(database, "up"); err != nil {
		slog.Error("Ошибка при применении миграций", "error", err)
		os.Exit(1)
	}

	repo := repositories.NewUserRepository(database)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	slog.Info("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":"+cfg.AppPort, router); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
