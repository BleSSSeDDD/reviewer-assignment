package service

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/GIT_USER_ID/GIT_REPO_ID/internal/storage"
)

// так как генерированный опенапи генератором код менять нельзя,
// то нельзя в мейне инициализировать бд,
//
//	 заранее неизвестно, какой хендлер будет дёрнут первым, а следовательно какая служебная
//	функция будет вызвана первой, поэтому приходится использовать глобальные переменные
//
// можно использовать инъекцию зависимостей  через структуру, но тут решил оставить переменные,
// зато получилось ленивое подключение бд
var (
	once      sync.Once
	db        *sql.DB
	dbOpenErr error
)

// нужен, чтобы не проверять в каждой функции что db != nil
func getDB() (*sql.DB, error) {
	once.Do(func() {
		db, dbOpenErr = storage.InitDB()
	})
	return db, dbOpenErr
}

// принимает строку с именем, проверяет нет ли такой команды и создает, либо возвращает ошибку
// так как db у нас глобальная переменная, то надо создавать локальную database,
// чтобы быть уверенным, что не произошло ошибки при подключении к бд
// длина строки не проверяется, потому что в бд уже есть ограничение
func CreateTeam(teamName string) error {
	exists, err := TeamExists(teamName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("team with that name already exists")
	}

	database, err := getDB()
	if err != nil {
		return err
	}

	_, err = database.Exec("INSERT INTO teams (name) VALUES ($1)", teamName)
	return err
}

// принимает строку с именем, проверяет есть ли такая команда в бд, возвращает тру/фолс и ошибку
func TeamExists(teamName string) (bool, error) {
	database, err := getDB()
	if err != nil {
		return false, err
	}

	var name string
	err = database.QueryRow("SELECT name FROM teams WHERE name = $1", teamName).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
