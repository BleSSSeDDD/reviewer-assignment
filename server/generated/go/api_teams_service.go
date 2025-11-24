package openapi

import (
	"context"
	"database/sql"
	"strings"

	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

// TeamsAPIService is a service that implements the logic for the TeamsAPIService
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
	transaction, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Response(500, nil), err
	}
	defer transaction.Rollback()

	err = storage.CreateTeam(ctx, transaction, team.TeamName)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return Response(409, ErrorResponse{
				Error: ErrorResponseError{
					Code:    "TEAM_EXISTS",
					Message: "team already exists",
				},
			}), nil
		}
		if strings.Contains(err.Error(), "value too long") {
			return Response(400, ErrorResponse{
				Error: ErrorResponseError{
					Code:    "VALUE_TOO_LONG",
					Message: "team name too long, maximum 100 characters",
				},
			}), nil
		}
		return Response(500, nil), err
	}

	for _, member := range team.Members {
		err = storage.CreateOrUpdateUser(ctx, transaction, member.UserId, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			if strings.Contains(err.Error(), "value too long") {
				return Response(400, ErrorResponse{
					Error: ErrorResponseError{
						Code:    "VALUE_TOO_LONG",
						Message: "user_id or username too long, maximum 100 characters",
					},
				}), nil
			}
			return Response(500, nil), err
		}

	}

	err = transaction.Commit()
	if err != nil {
		return Response(500, nil), err
	}

	return Response(201, TeamAddPost201Response{
		Team: team,
	}), nil
}

// TeamGetGet - Получить команду с участниками
func (s *TeamsAPIService) TeamGetGet(ctx context.Context, teamName string) (ImplResponse, error) {
	exists, err := storage.TeamExists(ctx, s.db, teamName)
	if err != nil {
		return Response(500, nil), err
	}
	if !exists {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "team not found",
			},
		}), nil
	}

	storageMembers, err := storage.GetTeamWithMembers(ctx, s.db, teamName)
	if err != nil {
		return Response(500, nil), err
	}

	var openapiMembers []TeamMember
	for _, storageMember := range storageMembers {
		openapiMembers = append(openapiMembers, TeamMember{
			UserId:   storageMember.UserId,
			Username: storageMember.Username,
			IsActive: storageMember.IsActive,
		})
	}

	return Response(200, Team{
		TeamName: teamName,
		Members:  openapiMembers,
	}), nil
}
