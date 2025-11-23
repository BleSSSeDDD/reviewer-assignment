package openapi

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
)

// TeamsAPIService is a service that implements the logic for the TeamsAPIServicer
// This service should implement the business logic for every endpoint for the TeamsAPI API.
// Include any external packages or services that will be required by this service.
type TeamsAPIService struct {
	db *sql.DB
}

// NewTeamsAPIService creates a default api service
func NewTeamsAPIService(db *sql.DB) *TeamsAPIService {
	return &TeamsAPIService{db: db}
}

// TeamAddPost - Создать команду с участниками (создаёт/обновляет пользователей)
func (s *TeamsAPIService) TeamAddPost(ctx context.Context, team Team) (ImplResponse, error) {
	// TODO - update TeamAddPost with the required logic for this service method.
	// Add api_teams_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(201, TeamAddPost201Response{}) or use other options such as http.Ok ...
	// return Response(201, TeamAddPost201Response{}), nil

	// TODO: Uncomment the next line to return response Response(400, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(400, ErrorResponse{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("TeamAddPost method not implemented")
}

// TeamGetGet - Получить команду с участниками
func (s *TeamsAPIService) TeamGetGet(ctx context.Context, teamName string) (ImplResponse, error) {
	// TODO - update TeamGetGet with the required logic for this service method.
	// Add api_teams_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Team{}) or use other options such as http.Ok ...
	// return Response(200, Team{}), nil

	// TODO: Uncomment the next line to return response Response(404, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(404, ErrorResponse{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("TeamGetGet method not implemented")
}
