package usecases

import (
	"context"
	"fmt"
	"pr-manager-service/internal/domain"
)

func (u *Usecases) GetStats(ctx context.Context) ([]domain.UserStats, []domain.PullRequestStats, error) {
	userStats, err := u.storage.GetUsersStats(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("GetUsersStats: %w", err)
	}

	pullRequestsStats, err := u.storage.GetPullRequestsStats(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("GetPullRequestsStats: %w", err)
	}

	return userStats, pullRequestsStats, nil
}
