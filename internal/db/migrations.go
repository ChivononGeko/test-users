package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func RunMigrations(db *sql.DB, direction string) error {
	slog.Info("Применение миграций...")

	migrationsDir := "./migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("ошибка при чтении папки с миграциями: %w", err)
	}

	var suffix string
	if direction == "up" {
		suffix = "_up.sql"
	} else if direction == "down" {
		suffix = "_down.sql"
	} else {
		return fmt.Errorf("неизвестное направление миграции: %s", direction)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), suffix) {
			continue
		}

		filePath := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("ошибка при чтении файла миграции %s: %w", filePath, err)
		}

		if _, err = db.Exec(string(content)); err != nil {
			return fmt.Errorf("ошибка при применении миграции %s: %w", filePath, err)
		}

		slog.Info("Миграция применена", "filename", file.Name())
	}

	slog.Info("Все миграции успешно применены.")
	return nil
}
