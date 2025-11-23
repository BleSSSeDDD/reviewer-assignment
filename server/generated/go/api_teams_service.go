package openapi

import (
	"context"
	"database/sql"
	"strings"

	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

type TeamsAPIService struct {
	db *sql.DB
}

func NewTeamsAPIService(db *sql.DB) *TeamsAPIService {
	return &TeamsAPIService{db: db}
}

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
		err = storage.CreateOrUpdateUser(ctx, transaction, member.UserId, member.Username, member.IsActive)
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

		err = storage.AddUserToTeam(ctx, transaction, member.UserId, team.TeamName)
		if err != nil {
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

func (s *TeamsAPIService) TeamGetGet(ctx context.Context, teamName string) (ImplResponse, error) {
	storageMembers, err := storage.GetTeamWithMembers(ctx, s.db, teamName)
	if err != nil {
		return Response(500, nil), err
	}

	if len(storageMembers) == 0 {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "team not found",
			},
		}), nil
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
