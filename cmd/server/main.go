package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	PullRequestsAPIService := server.NewPullRequestsAPIService()
	PullRequestsAPIController := server.NewPullRequestsAPIController(PullRequestsAPIService)

	TeamsAPIService := server.NewTeamsAPIService()
	TeamsAPIController := server.NewTeamsAPIController(TeamsAPIService)

	UsersAPIService := server.NewUsersAPIService()
	UsersAPIController := server.NewUsersAPIController(UsersAPIService)

	router := server.NewRouter(PullRequestsAPIController, TeamsAPIController, UsersAPIController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
