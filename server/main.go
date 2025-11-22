package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	openapi "github.com/BleSSSeDDD/reviewer-assignment/server/generated/go"
)

func setupRouter() http.Handler {
	PullRequestsAPIService := openapi.NewPullRequestsAPIService()
	PullRequestsAPIController := openapi.NewPullRequestsAPIController(PullRequestsAPIService)

	TeamsAPIService := openapi.NewTeamsAPIService()
	TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)

	UsersAPIService := openapi.NewUsersAPIService()
	UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)

	return openapi.NewRouter(PullRequestsAPIController, TeamsAPIController, UsersAPIController)
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stop := make(chan os.Signal, 1)                    // канал для грейсфул шатдауна
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM) // сигтерм для докера

	serverError := make(chan error, 1) // канал для завершения в случае если сервер перестанет работать

	log.Printf("Server starting")

	router := setupRouter()
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

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}

	log.Println("Server stopped")
}
