package openapi

import (
	"context"
	"database/sql"

	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

// UsersAPIService is a service that implements the logic for the UsersAPIServicer
// This service should implement the business logic for every endpoint for the UsersAPI API.
// Include any external packages or services that will be required by this service.
type UsersAPIService struct {
	db *sql.DB
}

// NewUsersAPIService creates a default api service
func NewUsersAPIService(db *sql.DB) *UsersAPIService {
	return &UsersAPIService{db: db}
}

// получает юзера из бд, если его нет то 404, если он есть, то обновляем is_active,
// потом получаем обновленную версию и отдаем в хендлер, причем надо конвертировать,
// потому что в сторедже своя структура под него
func (s *UsersAPIService) UsersSetIsActivePost(ctx context.Context, req UsersSetIsActivePostRequest) (ImplResponse, error) {
	storageUser, err := storage.GetUser(ctx, s.db, req.UserId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
		}), nil
	}

	err = storage.CreateOrUpdateUser(ctx, s.db, req.UserId, storageUser.Username, req.IsActive)
	if err != nil {
		return Response(500, nil), err
	}

	updatedStorageUser, _ := storage.GetUser(ctx, s.db, req.UserId)

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

// проверяем что юзер вообще существует, получает слайс структур пулреквеста на уровне бд, где этот юзер ревьювер,
// потом их конвертирует в нормальную структуру и отдает наверх
func (s *UsersAPIService) UsersGetReviewGet(ctx context.Context, userId string) (ImplResponse, error) {
	_, err := storage.GetUser(ctx, s.db, userId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
		}), nil
	}

	storagePRs, err := storage.GetUserReviewPRs(ctx, s.db, userId)
	if err != nil {
		return Response(500, nil), err
	}

	var openapiPRs []PullRequestShort
	for _, storagePR := range storagePRs {
		openapiPRs = append(openapiPRs, PullRequestShort{
			PullRequestId:   storagePR.PullRequestId,
			PullRequestName: storagePR.PullRequestName,
			AuthorId:        storagePR.AuthorId,
			Status:          storagePR.Status,
		})
	}

	return Response(200, UsersGetReviewGet200Response{
		UserId:       userId,
		PullRequests: openapiPRs,
	}), nil
}
