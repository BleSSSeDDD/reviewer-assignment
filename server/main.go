package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	openapi "github.com/BleSSSeDDD/reviewer-assignment/server/generated/go"
	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

func setupRouter(db *sql.DB) http.Handler {
	PullRequestsAPIService := openapi.NewPullRequestsAPIService(db)
	PullRequestsAPIController := openapi.NewPullRequestsAPIController(PullRequestsAPIService)

	TeamsAPIService := openapi.NewTeamsAPIService(db)
	TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)

	UsersAPIService := openapi.NewUsersAPIService(db)
	UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)

	return openapi.NewRouter(PullRequestsAPIController, TeamsAPIController, UsersAPIController)
}

func main() {

	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = storage.InitDB()
		if err != nil {
			log.Printf("Server could not connect to db, retrying...")
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Printf("Server could not connect to db")
		return
	}

	stop := make(chan os.Signal, 1)                    // канал для грейсфул шатдауна
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM) // сигтерм для докера

	serverError := make(chan error, 1) // канал для завершения в случае если сервер перестанет работать

	log.Printf("Server starting")

	router := setupRouter(db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// запуск в горутине чтобы можно было дальше в мейне ждать сигналы для завершения программы
	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil {
			serverError <- err
		}
	}()

	select {
	case <-stop:
		log.Println("Received shutdown signal")
	case err := <-serverError:
		log.Printf("Server error: %v\n", err)
	}

	log.Println("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		server.Shutdown(context.Background())
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}

	log.Println("Server stopped")
}
