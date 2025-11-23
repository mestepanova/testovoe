package usecases

//go:generate go tool mockgen -source=./contract.go -destination=./contract_mock.go -package=usecases

import (
	"context"
	"pr-manager-service/internal/domain"
)

type Storage interface {
	CreateTeam(ctx context.Context, request domain.CreateTeamRequest, teamID string) error
	GetTeamByName(ctx context.Context, teamName string) (domain.Team, error)
	GetTeamFullByName(ctx context.Context, teamName string) (domain.Team, []domain.User, error)
	GetActiveColleagues(ctx context.Context, userID string) ([]domain.User, error)

	GetPullRequestsByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error)
	GetPullRequestByID(ctx context.Context, prID string) (domain.PullRequest, error)
	CreatePullRequest(
		ctx context.Context,
		request domain.CreatePullRequestRequest,
		reviewersIDs []string,
	) (pr domain.PullRequest, err error)
	UpdatePullRequestStatus(ctx context.Context, prID string, newStatus domain.PullRequestStatus) error
	UpdatePullRequestReviewersIDs(ctx context.Context, prID string, reviewersIDs []string) error

	CreateUsers(ctx context.Context, requests []domain.CreateUserRequest, teamID string) error
	UpdateUserStatus(ctx context.Context, userID string, isActive bool) error
	GetUserFull(ctx context.Context, userID string) (domain.User, domain.Team, error)
	GetUserShort(ctx context.Context, userID string) (domain.User, error)

	PullRequestStatsCreate(ctx context.Context, pullRequestID string, assignmentsCount int) error
	UserStatsCreateBatch(ctx context.Context, userIDs []string) error
	UserAssignmentsIncrementBatch(ctx context.Context, userIDs []string) error
	UserStatusChangesIncrementBatch(ctx context.Context, userIDs []string) error
	PullRequestAssignmentsIncrement(ctx context.Context, pullRequestID string) error
	GetUsersStats(ctx context.Context) (userStats []domain.UserStats, err error)
	GetPullRequestsStats(ctx context.Context) (pullRequestsStats []domain.PullRequestStats, err error)

	UnitOfWork(ctx context.Context, do func(s Storage) error) error
}
