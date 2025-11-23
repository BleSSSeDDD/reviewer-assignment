package openapi

import (
	"context"
	"database/sql"
	"strings"

	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

type UsersAPIService struct {
	db *sql.DB
}

func NewUsersAPIService(db *sql.DB) *UsersAPIService {
	return &UsersAPIService{db: db}
}

func (s *UsersAPIService) UsersSetIsActivePost(ctx context.Context, req UsersSetIsActivePostRequest) (ImplResponse, error) {
	transaction, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Response(500, nil), err
	}
	defer transaction.Rollback()

	storageUser, err := storage.GetUser(ctx, transaction, req.UserId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
		}), nil
	}

	err = storage.CreateOrUpdateUser(ctx, transaction, req.UserId, storageUser.Username, req.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "value too long") {
			return Response(400, ErrorResponse{
				Error: ErrorResponseError{
					Code:    "VALUE_TOO_LONG",
					Message: "user_id too long, maximum 100 characters",
				},
			}), nil
		}
		return Response(500, nil), err
	}

	updatedStorageUser, err := storage.GetUser(ctx, transaction, req.UserId)
	if err != nil {
		return Response(500, nil), err
	}

	err = transaction.Commit()
	if err != nil {
		return Response(500, nil), err
	}

	updatedUser := User{
		UserId:   updatedStorageUser.UserId,
		Username: updatedStorageUser.Username,
		IsActive: updatedStorageUser.IsActive,
		TeamName: updatedStorageUser.TeamName,
	}

	return Response(200, UsersSetIsActivePost200Response{
		User: updatedUser,
	}), nil
}

func (s *UsersAPIService) UsersGetReviewGet(ctx context.Context, userId string) (ImplResponse, error) {
	storagePullRequests, err := storage.GetUserReviewPRs(ctx, s.db, userId)
	if err != nil {
		return Response(500, nil), err
	}

	var openapiPullRequests []PullRequestShort
	for _, storagePullRequest := range storagePullRequests {
		openapiPullRequests = append(openapiPullRequests, PullRequestShort{
			PullRequestId:   storagePullRequest.PullRequestId,
			PullRequestName: storagePullRequest.PullRequestName,
			AuthorId:        storagePullRequest.AuthorId,
			Status:          storagePullRequest.Status,
		})
	}

	return Response(200, UsersGetReviewGet200Response{
		UserId:       userId,
		PullRequests: openapiPullRequests,
	}), nil
}
