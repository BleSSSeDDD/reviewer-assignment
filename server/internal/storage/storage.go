package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// структура юзера, поля такие же как в бд, чтобы не тащить его структуру со слоя выше
type StorageUser struct {
	UserId   string
	Username string
	IsActive bool
	TeamName string
}

// такая же структура для пул реквеста
type StoragePullRequest struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          string
}

// инициализируем бд, если всё ок, то возвращается указатель на бд и nil,
// иначе nil и ошибка и надо завершать программу
func InitDB() (db *sql.DB, err error) {
	connect := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), "5432", os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	log.Printf("Connecting to DB: host=%s, dbname=%s, user=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"))

	db, err = sql.Open("postgres", connect)
	if err != nil {
		log.Printf("SQL Open error: %v", err)
		return nil, fmt.Errorf("cant open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("DB Ping error: %v", err)
		return nil, fmt.Errorf("cant ping db: %v", err)
	}

	log.Printf("Successfully connected to DB")
	return db, nil
}

// получает юзера по айди, если его нет то вернет ошибку
// team_name берется из members_of_teams, юзер не может быть не в команде
func GetUser(ctx context.Context, db *sql.DB, userID string) (StorageUser, error) {
	var user StorageUser
	err := db.QueryRowContext(ctx, `
        SELECT u.user_id, u.user_name, u.is_active, m.team_name 
        FROM users u
        LEFT JOIN members_of_teams m ON u.user_id = m.user_id
        WHERE u.user_id = $1
    `, userID).Scan(&user.UserId, &user.Username, &user.IsActive, &user.TeamName)
	return user, err
}

// создает юзера если его нет, либо обновляет если уже есть и нам надо поле isActive поменять
func CreateOrUpdateUser(ctx context.Context, db *sql.DB, userID, userName string, isActive bool) error {
	_, err := db.ExecContext(ctx, `
        INSERT INTO users (user_id, user_name, is_active) 
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE SET is_active = EXCLUDED.is_active
    `, userID, userName, isActive)
	return err
}

// ищет в бд юзера, потом его пул реквсты и возвращает их слайс и ошибку если что-то пошло не так
func GetUserReviewPRs(ctx context.Context, db *sql.DB, userID string) ([]StoragePullRequest, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT pr.request_id, pr.title, pr.user_id, pr.status
		FROM pull_requests pr
		JOIN pull_requests_reviewers prr ON pr.request_id = prr.request_id
		WHERE prr.reviewer_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var UserPullRequests []StoragePullRequest
	for rows.Next() {
		var pullRequest StoragePullRequest
		if err := rows.Scan(&pullRequest.PullRequestId, &pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.Status); err != nil {
			return nil, err
		}
		UserPullRequests = append(UserPullRequests, pullRequest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UserPullRequests, nil
}
