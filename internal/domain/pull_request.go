package domain

import (
	"pr-manager-service/internal/generated/api"
	"time"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorUserID      string
	ReviewersUsersIDs []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
	Status            PullRequestStatus
}

type PullRequestStatus uint8

const (
	StatusOpen   PullRequestStatus = 0
	StatusMerged PullRequestStatus = 1
)

func ConvertPullRequestStatusToDomain(status api.PullRequestStatus) PullRequestStatus {
	switch status {
	case api.OPEN:
		return StatusOpen
	case api.MERGED:
		return StatusMerged
	}
	return StatusOpen
}

func ConvertPullRequestStatusToApi(status PullRequestStatus) api.PullRequestStatus {
	switch status {
	case StatusOpen:
		return api.OPEN
	case StatusMerged:
		return api.MERGED
	}
	return api.OPEN
}

func ConvertPullRequest(pr PullRequest) api.PullRequest {
	return api.PullRequest{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorId:          pr.AuthorUserID,
		AssignedReviewers: pr.ReviewersUsersIDs,
		Status:            ConvertPullRequestStatusToApi(pr.Status),
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

type CreatePullRequestRequest struct {
	ID           string `json:"pull_request_id"   validate:"required,min=1,max=36"`
	Name         string `json:"pull_request_name" validate:"required,min=2,max=50"`
	AuthorUserID string `json:"author_id"         validate:"required,min=1,max=36"`
}
