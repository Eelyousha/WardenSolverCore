package main

import (
	v1 "handlerv1"
	"log"
	"net/http"
)

// Database подключение к БД (глобальное)
var db *Database

func main() {
	// Инициализация подключения к БД
	database, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	db = database

	// Запуск миграций
	err = RunMigrations(database.conn, "migrations")
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Регистрация обработчиков
	http.HandleFunc("/", v1.HandleRoot)
	http.HandleFunc("/api/post", v1.HandlePostBlacklist(db))
	http.HandleFunc("/api/blacklist", v1.HandleGetBlacklist(db))
	http.HandleFunc("/api/check", v1.HandleCheckBlacklist(db))

	// Запуск сервера
	port := ":8080"
	log.Printf("Server started at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
