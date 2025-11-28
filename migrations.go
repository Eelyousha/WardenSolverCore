package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations выполняет все миграции из папки migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Получить список всех .sql файлов
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// Отсортировать файлы по имени (чтобы выполнять в правильном порядке)
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("No migration files found")
		return nil
	}

	// Выполнить каждый файл миграции
	for _, sqlFile := range sqlFiles {
		filePath := filepath.Join(migrationsPath, sqlFile)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", sqlFile, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", sqlFile, err)
		}

		log.Printf("✓ Migration applied: %s\n", sqlFile)
	}

	return nil
}
