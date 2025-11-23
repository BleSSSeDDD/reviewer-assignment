package openapi

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

// PullRequestsAPIService is a service that implements the logic for the PullRequestsAPIServicer
// This service should implement the business logic for every endpoint for the PullRequestsAPI API.
// Include any external packages or services that will be required by this service.
type PullRequestsAPIService struct {
	db *sql.DB
}

// NewPullRequestsAPIService creates a default api service
func NewPullRequestsAPIService(db *sql.DB) *PullRequestsAPIService {
	return &PullRequestsAPIService{db: db}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PullRequestCreatePost - Создать PullRequest и автоматически назначить до 2 ревьюверов из команды автора
func (s *PullRequestsAPIService) PullRequestCreatePost(ctx context.Context, pullRequestCreatePostRequest PullRequestCreatePostRequest) (ImplResponse, error) {
	transaction, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Response(500, nil), err
	}
	defer transaction.Rollback()

	authorUser, err := storage.GetUser(ctx, transaction, pullRequestCreatePostRequest.AuthorId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "author not found",
			},
		}), nil
	}

	teamMembers, err := storage.GetTeamWithMembers(ctx, s.db, authorUser.TeamName)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "author team not found",
			},
		}), nil
	}

	var availableReviewers []storage.StorageUser
	for _, member := range teamMembers {
		if member.UserId != authorUser.UserId && member.IsActive {
			availableReviewers = append(availableReviewers, member)
		}
	}

	// чтобы случайно выбрать из юзеров, которые не автор пул реквеста делаем слайс,
	// потом его перемешиваем и берем двух первых, либо если там один всего то его
	var selectedReviewers []storage.StorageUser
	if len(availableReviewers) > 0 {
		rand.Shuffle(len(availableReviewers), func(i, j int) {
			availableReviewers[i], availableReviewers[j] = availableReviewers[j], availableReviewers[i]
		})

		count := min(2, len(availableReviewers))
		selectedReviewers = availableReviewers[:count]
	}

	err = storage.CreatePullRequest(ctx, transaction,
		pullRequestCreatePostRequest.PullRequestId,
		pullRequestCreatePostRequest.PullRequestName,
		pullRequestCreatePostRequest.AuthorId,
	)
	if err != nil {
		return Response(409, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "PullRequest_EXISTS",
				Message: "PullRequest id already exists",
			},
		}), nil
	}

	var assignedReviewerIds []string
	for _, reviewer := range selectedReviewers {
		err = storage.AddReviewerToPR(ctx, transaction, pullRequestCreatePostRequest.PullRequestId, reviewer.UserId)
		if err != nil {
			return Response(500, nil), err
		}
		assignedReviewerIds = append(assignedReviewerIds, reviewer.UserId)
	}

	createdPullRequest, err := storage.GetPullRequest(ctx, transaction, pullRequestCreatePostRequest.PullRequestId)
	if err != nil {
		return Response(500, nil), err
	}

	err = transaction.Commit()
	if err != nil {
		return Response(500, nil), err
	}

	openapiPullRequest := PullRequest{
		PullRequestId:     createdPullRequest.PullRequestId,
		PullRequestName:   createdPullRequest.PullRequestName,
		AuthorId:          createdPullRequest.AuthorId,
		Status:            createdPullRequest.Status,
		AssignedReviewers: assignedReviewerIds,
		CreatedAt:         createdPullRequest.CreatedAt,
		MergedAt:          createdPullRequest.MergedAt,
	}

	return Response(201, PullRequestCreatePost201Response{
		Pr: openapiPullRequest,
	}), nil
}

// PullRequestMergePost - Пометить PullRequest как MERGED (идемпотентная операция)
func (s *PullRequestsAPIService) PullRequestMergePost(ctx context.Context, pullRequestMergePostRequest PullRequestMergePostRequest) (ImplResponse, error) {
	transaction, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Response(500, nil), err
	}
	defer transaction.Rollback()

	storagePullRequest, err := storage.GetPullRequest(ctx, transaction, pullRequestMergePostRequest.PullRequestId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "PullRequest not found",
			},
		}), nil
	}

	if storagePullRequest.Status == "MERGED" {
		reviewers, err := storage.GetPRReviewers(ctx, s.db, storagePullRequest.PullRequestId)
		if err != nil {
			return Response(500, nil), err
		}

		openapiPullRequest := PullRequest{
			PullRequestId:     storagePullRequest.PullRequestId,
			PullRequestName:   storagePullRequest.PullRequestName,
			AuthorId:          storagePullRequest.AuthorId,
			Status:            storagePullRequest.Status,
			AssignedReviewers: reviewers,
			CreatedAt:         storagePullRequest.CreatedAt,
			MergedAt:          storagePullRequest.MergedAt,
		}

		return Response(200, PullRequestCreatePost201Response{
			Pr: openapiPullRequest,
		}), nil
	}

	err = storage.UpdatePRStatus(ctx, transaction, pullRequestMergePostRequest.PullRequestId, "MERGED")
	if err != nil {
		return Response(500, nil), err
	}

	mergedPullRequest, err := storage.GetPullRequest(ctx, transaction, pullRequestMergePostRequest.PullRequestId)
	if err != nil {
		return Response(500, nil), err
	}

	reviewers, err := storage.GetPRReviewers(ctx, s.db, mergedPullRequest.PullRequestId)
	if err != nil {
		return Response(500, nil), err
	}

	err = transaction.Commit()
	if err != nil {
		return Response(500, nil), err
	}

	openapiPullRequest := PullRequest{
		PullRequestId:     mergedPullRequest.PullRequestId,
		PullRequestName:   mergedPullRequest.PullRequestName,
		AuthorId:          mergedPullRequest.AuthorId,
		Status:            mergedPullRequest.Status,
		AssignedReviewers: reviewers,
		CreatedAt:         mergedPullRequest.CreatedAt,
		MergedAt:          mergedPullRequest.MergedAt,
	}

	return Response(200, PullRequestCreatePost201Response{
		Pr: openapiPullRequest,
	}), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PullRequestReassignPost - Переназначить конкретного ревьювера на другого из его команды
func (s *PullRequestsAPIService) PullRequestReassignPost(ctx context.Context, pullRequestReassignPostRequest PullRequestReassignPostRequest) (ImplResponse, error) {
	transaction, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Response(500, nil), err
	}
	defer transaction.Rollback()

	storagePullRequest, err := storage.GetPullRequest(ctx, transaction, pullRequestReassignPostRequest.PullRequestId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "PullRequest not found",
			},
		}), nil
	}

	if storagePullRequest.Status == "MERGED" {
		return Response(409, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "PullRequest_MERGED",
				Message: "cannot reassign on merged PullRequest",
			},
		}), nil
	}

	currentReviewers, err := storage.GetPRReviewers(ctx, s.db, storagePullRequest.PullRequestId)
	if err != nil {
		return Response(500, nil), err
	}

	isAssigned := false
	for _, reviewer := range currentReviewers {
		if reviewer == pullRequestReassignPostRequest.OldUserId {
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		return Response(409, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_ASSIGNED",
				Message: "reviewer is not assigned to this PullRequest",
			},
		}), nil
	}

	// берем команду старого ревьюера и получаем активных пользователей из его команды
	oldReviewerUser, err := storage.GetUser(ctx, transaction, pullRequestReassignPostRequest.OldUserId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "old reviewer not found",
			},
		}), nil
	}

	teamMembers, err := storage.GetTeamWithMembers(ctx, s.db, oldReviewerUser.TeamName)
	if err != nil {
		return Response(500, nil), err
	}

	var availableCandidates []storage.StorageUser
	for _, member := range teamMembers {
		if member.IsActive &&
			member.UserId != pullRequestReassignPostRequest.OldUserId &&
			member.UserId != storagePullRequest.AuthorId &&
			!contains(currentReviewers, member.UserId) {
			availableCandidates = append(availableCandidates, member)
		}
	}

	if len(availableCandidates) == 0 {
		return Response(409, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NO_CANDIDATE",
				Message: "no active replacement candidate in team",
			},
		}), nil
	}

	// он всего один, перетасовывать слайс не надо
	newReviewerUser := availableCandidates[rand.Intn(len(availableCandidates))]
	err = storage.ReplacePRReviewer(ctx, transaction,
		pullRequestReassignPostRequest.PullRequestId,
		pullRequestReassignPostRequest.OldUserId,
		newReviewerUser.UserId,
	)
	if err != nil {
		return Response(500, nil), err
	}

	updatedReviewers, err := storage.GetPRReviewers(ctx, s.db, storagePullRequest.PullRequestId)
	if err != nil {
		return Response(500, nil), err
	}

	err = transaction.Commit()
	if err != nil {
		return Response(500, nil), err
	}

	openapiPullRequest := PullRequest{
		PullRequestId:     storagePullRequest.PullRequestId,
		PullRequestName:   storagePullRequest.PullRequestName,
		AuthorId:          storagePullRequest.AuthorId,
		Status:            storagePullRequest.Status,
		AssignedReviewers: updatedReviewers,
		CreatedAt:         storagePullRequest.CreatedAt,
		MergedAt:          storagePullRequest.MergedAt,
	}

	return Response(200, PullRequestReassignPost200Response{
		Pr:         openapiPullRequest,
		ReplacedBy: newReviewerUser.UserId,
	}), nil
}
