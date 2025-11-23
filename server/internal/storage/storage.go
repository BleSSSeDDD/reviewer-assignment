package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
	CreatedAt       *time.Time
	MergedAt        *time.Time
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
// транзакция потому что мы потом будем его менять и надо чтобы он не изменился пока мы работаем
func GetUser(ctx context.Context, transaction *sql.Tx, userId string) (StorageUser, error) {
	var user StorageUser
	err := transaction.QueryRowContext(ctx, `
        SELECT users.user_id, users.user_name, users.is_active, members_of_teams.team_name 
        FROM users 
        LEFT JOIN members_of_teams ON users.user_id = members_of_teams.user_id
        WHERE users.user_id = $1
    `, userId).Scan(&user.UserId, &user.Username, &user.IsActive, &user.TeamName)
	return user, err
}

// создает юзера если его нет, либо обновляет если уже есть и нам надо поле isActive поменять
// транзакция потому что связавно с гетюзером
func CreateOrUpdateUser(ctx context.Context, transaction *sql.Tx, userId, userName string, isActive bool) error {
	_, err := transaction.ExecContext(ctx, `
        INSERT INTO users (user_id, user_name, is_active) 
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE SET is_active = EXCLUDED.is_active
    `, userId, userName, isActive)
	return err
}

// ищет в бд юзера, потом его пул реквсты и возвращает их слайс и ошибку если что-то пошло не так
// не транзакция потому что мы просто читаем и ничего не меняем
func GetUserReviewPRs(ctx context.Context, db *sql.DB, userId string) ([]StoragePullRequest, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT pull_requests.request_id, pull_requests.title, pull_requests.user_id, pull_requests.status, pull_requests.created_at, pull_requests.merged_at
		FROM pull_requests 
		JOIN pull_requests_reviewers ON pull_requests.request_id = pull_requests_reviewers.request_id
		WHERE pull_requests_reviewers.reviewer_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var UserPullRequests []StoragePullRequest
	for rows.Next() {
		var pullRequest StoragePullRequest
		if err := rows.Scan(&pullRequest.PullRequestId, &pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.Status, &pullRequest.CreatedAt, &pullRequest.MergedAt); err != nil {
			return nil, err
		}
		UserPullRequests = append(UserPullRequests, pullRequest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UserPullRequests, nil
}

// создает команду если её нет, если есть возвращает ошибку
// транзакция потому что потом сюда надо добавить юзеров
func CreateTeam(ctx context.Context, transaction *sql.Tx, teamName string) error {
	_, err := transaction.ExecContext(ctx, `
        INSERT INTO teams (team_name) 
        VALUES ($1)
    `, teamName)
	return err
}

// добавляет юзера в команду, если уже есть то ошибка
// транзакция потому что это должно быть атомарным с созданием команды и юзеров
func AddUserToTeam(ctx context.Context, transaction *sql.Tx, userID, teamName string) error {
	_, err := transaction.ExecContext(ctx, `
        INSERT INTO members_of_teams (user_id, team_name) 
        VALUES ($1, $2)
    `, userID, teamName)
	return err
}

// получает всех юзеров команды, если команды нет то пустой слайс
// не транзакция потому что мы просто читаем данные и нам не важно если они потом изменятся
func GetTeamWithMembers(ctx context.Context, db *sql.DB, teamName string) ([]StorageUser, error) {
	rows, err := db.QueryContext(ctx, `
        SELECT users.user_id, users.user_name, users.is_active, members_of_teams.team_name
        FROM members_of_teams 
        JOIN users ON members_of_teams.user_id = users.user_id
        WHERE members_of_teams.team_name = $1
    `, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []StorageUser
	for rows.Next() {
		var user StorageUser
		if err := rows.Scan(&user.UserId, &user.Username, &user.IsActive, &user.TeamName); err != nil {
			return nil, err
		}
		members = append(members, user)
	}
	return members, nil
}

// создает пул реквест если его нет, если есть то ошибка
// транзакция потому что потом сюда надо добавить ревьюеров
func CreatePullRequest(ctx context.Context, transaction *sql.Tx, pullRequestId, title, authorId string) error {
	_, err := transaction.ExecContext(ctx, `
        INSERT INTO pull_requests (request_id, title, user_id, status) 
        VALUES ($1, $2, $3, 'OPEN')
    `, pullRequestId, title, authorId)
	return err
}

// добавляет ревьюера к пул реквесту, если уже есть то ошибка
// транзакция потому что это должно быть атомарным с созданием пул реквеста
func AddReviewerToPR(ctx context.Context, transaction *sql.Tx, pullRequestId, reviewerId string) error {
	_, err := transaction.ExecContext(ctx, `
        INSERT INTO pull_requests_reviewers (request_id, reviewer_id) 
        VALUES ($1, $2)
    `, pullRequestId, reviewerId)
	return err
}

// получает пул реквест по айди, если нет то ошибка
// транзакция потому что мы потом будем его менять
func GetPullRequest(ctx context.Context, transaction *sql.Tx, pullRequestId string) (StoragePullRequest, error) {
	var pullRequest StoragePullRequest
	err := transaction.QueryRowContext(ctx, `
        SELECT request_id, title, user_id, status, created_at, merged_at
        FROM pull_requests 
        WHERE request_id = $1
    `, pullRequestId).Scan(&pullRequest.PullRequestId, &pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.Status, &pullRequest.CreatedAt, &pullRequest.MergedAt)
	return pullRequest, err
}

// обновляет статус пул реквеста, если его нет то ошибка
// транзакция потому что это изменение данных
func UpdatePRStatus(ctx context.Context, transaction *sql.Tx, pullRequestId, status string) error {
	if status == "MERGED" {
		_, err := transaction.ExecContext(ctx, `
            UPDATE pull_requests 
            SET status = $1, merged_at = NOW()
            WHERE request_id = $2
        `, status, pullRequestId)
		return err
	}

	_, err := transaction.ExecContext(ctx, `
        UPDATE pull_requests 
        SET status = $1 
        WHERE request_id = $2
    `, status, pullRequestId)
	return err
}

// получает всех ревьюеров пул реквеста, если пул реквеста нет то пустой слайс
func GetPRReviewers(ctx context.Context, db *sql.DB, pullRequestId string) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
        SELECT reviewer_id 
        FROM pull_requests_reviewers 
        WHERE request_id = $1
    `, pullRequestId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerId string
		if err := rows.Scan(&reviewerId); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewerId)
	}
	return reviewers, nil
}

// удаляет ревьюера из пул реквеста и добавляет нового
// транзакция потому что это два изменения которые должны быть атомарными
func ReplacePRReviewer(ctx context.Context, transaction *sql.Tx, pullRequestId, oldReviewerId, newReviewerId string) error {
	_, err := transaction.ExecContext(ctx, `
        DELETE FROM pull_requests_reviewers 
        WHERE request_id = $1 AND reviewer_id = $2
    `, pullRequestId, oldReviewerId)
	if err != nil {
		return err
	}

	_, err = transaction.ExecContext(ctx, `
        INSERT INTO pull_requests_reviewers (request_id, reviewer_id) 
        VALUES ($1, $2)
    `, pullRequestId, newReviewerId)
	return err
}
