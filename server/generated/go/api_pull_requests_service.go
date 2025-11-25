package openapi

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"

	"github.com/BleSSSeDDD/reviewer-assignment/server/internal/storage"
)

type PullRequestsAPIService struct {
	db *sql.DB
}

func NewPullRequestsAPIService(db *sql.DB) *PullRequestsAPIService {
	return &PullRequestsAPIService{db: db}
}

// PullRequestCreatePost - Создать PR и автоматически назначить до 2 ревьюверов из команды автора
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

	teamMembers, err := storage.GetTeamWithMembers(ctx, transaction, authorUser.TeamName)
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
		if strings.Contains(err.Error(), "duplicate key") {
			return Response(409, ErrorResponse{
				Error: ErrorResponseError{
					Code:    "PR_EXISTS",
					Message: "PullRequest id already exists",
				},
			}), nil
		}
		if strings.Contains(err.Error(), "value too long") {
			return Response(400, ErrorResponse{
				Error: ErrorResponseError{
					Code:    "VALUE_TOO_LONG",
					Message: "PR id or title too long, check maximum length limits",
				},
			}), nil
		}
		return Response(500, nil), err
	}

	assignedReviewerIds := make([]string, 0)
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

// PullRequestMergePost - Пометить PR как MERGED (идемпотентная операция)
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
		reviewers, err := storage.GetPRReviewers(ctx, transaction, storagePullRequest.PullRequestId)
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

	reviewers, err := storage.GetPRReviewers(ctx, transaction, mergedPullRequest.PullRequestId)
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
				Code:    "PR_MERGED",
				Message: "cannot reassign on merged PullRequest",
			},
		}), nil
	}

	currentReviewers, err := storage.GetPRReviewers(ctx, transaction, storagePullRequest.PullRequestId)
	if err != nil {
		return Response(500, nil), err
	}

	if !contains(currentReviewers, pullRequestReassignPostRequest.OldUserId) {
		return Response(409, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_ASSIGNED",
				Message: "reviewer is not assigned to this PullRequest",
			},
		}), nil
	}

	oldReviewerUser, err := storage.GetUser(ctx, transaction, pullRequestReassignPostRequest.OldUserId)
	if err != nil {
		return Response(404, ErrorResponse{
			Error: ErrorResponseError{
				Code:    "NOT_FOUND",
				Message: "old reviewer not found",
			},
		}), nil
	}

	teamMembers, err := storage.GetTeamWithMembers(ctx, transaction, oldReviewerUser.TeamName)
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

	newReviewerUser := availableCandidates[rand.Intn(len(availableCandidates))]
	err = storage.ReplacePRReviewer(ctx, transaction,
		pullRequestReassignPostRequest.PullRequestId,
		pullRequestReassignPostRequest.OldUserId,
		newReviewerUser.UserId,
	)
	if err != nil {
		return Response(500, nil), err
	}

	updatedReviewers, err := storage.GetPRReviewers(ctx, transaction, storagePullRequest.PullRequestId)
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
