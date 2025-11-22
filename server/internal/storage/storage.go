package storage

import (
	"database/sql"
	"fmt"
	"os"
)

// так как генерированный опенапи генератором код менять нельзя, то нельзя в мейне инициализировать бд,
// заранее неизвестно, какой хендлер будет дёрнут первым, а следовательно какая служебная
// функция будет вызвана первой, поэтому приходится использовать глобальные переменные
// можно использовать инъекцию зависимостей  через структуру, но тут решил оставить переменные,
// зато получилось ленивое подключение бд :)
var (
	db        *sql.DB
	dbOpenErr error
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

// принимает строку с именем и создает, либо возвращает ошибку от бд
// так как db у нас глобальная переменная, то надо создавать локальную database,
// чтобы быть уверенным, что не произошло ошибки при подключении к бд
// длина строки не проверяется, потому что в бд уже есть ограничение
func CreateTeam(teamName string) error {
	if db == nil {
		db, dbOpenErr = InitDB()
		if dbOpenErr != nil {
			return dbOpenErr
		}
	}

	_, err := db.Exec("INSERT INTO teams (name) VALUES ($1)", teamName)
	return err
}

// принимает айди пользователя, имя и статус активности,
// если что-то пошло не так при создании(например он уже был в бд), то вернет ошибку от бд
func CreateUser(user_id, user_name string, is_active bool) error {
	if db == nil {
		db, dbOpenErr = InitDB()
		if dbOpenErr != nil {
			return dbOpenErr
		}
	}

	_, err := db.Exec(
		"INSERT INTO users (user_id, user_name, is_active) VALUES ($1, $2, $3)",
		user_id, user_name, is_active,
	)
	return err
}
