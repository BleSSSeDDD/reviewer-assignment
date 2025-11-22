package storage

import (
	"database/sql"
	"fmt"
	"os"
)

// инициализируем бд, если всё ок, то возвращается указатель на бд и nil,
// иначе nil и ошибка и надо завершать программу
func InitDB() (db *sql.DB, err error) {
	connect := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"), "5432", os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("cant open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cant ping db: %v", err)
	}

	return db, nil
}
