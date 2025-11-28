package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Database представляет подключение к PostgreSQL
type Database struct {
	conn *sql.DB
}

// NewDatabase создаёт новое подключение к БД
func NewDatabase() (*Database, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "warden"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{conn: db}, nil
}

// Close закрывает подключение к БД
func (d *Database) Close() error {
	return d.conn.Close()
}

// AddToBlacklist добавляет пользователя в чёрный список
func (d *Database) AddToBlacklist(socialNetwork, userID, nickname string) (int, error) {
	var id int
	query := `INSERT INTO blacklist (social_network, user_id, nickname)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (social_network, user_id, nickname) DO NOTHING
			 RETURNING id;`

	err := d.conn.QueryRow(query, socialNetwork, userID, nickname).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Пользователь уже в чёрном списке
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

// GetBlacklist получает список всех забаненных пользователей
func (d *Database) GetBlacklist(socialNetwork, userID string) ([]map[string]interface{}, error) {
	query := `SELECT nickname, created_at
			 FROM blacklist
			 WHERE social_network = $1 AND user_id = $2
			 ORDER BY created_at DESC;`

	rows, err := d.conn.Query(query, socialNetwork, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var nickname string
		var createdAt time.Time

		err := rows.Scan(&nickname, &createdAt)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"nickname":   nickname,
			"created_at": createdAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// CheckBlacklist проверяет наличие пользователя в чёрном списке
func (d *Database) CheckBlacklist(socialNetwork, userID, nickname string) (bool, error) {
	query := `SELECT EXISTS(
				 SELECT 1 FROM blacklist
				 WHERE social_network = $1 AND user_id = $2 AND nickname = $3
			 ) AS is_banned;`

	var isBanned bool
	err := d.conn.QueryRow(query, socialNetwork, userID, nickname).Scan(&isBanned)
	if err != nil {
		return false, err
	}

	return isBanned, nil
}
