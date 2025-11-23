package usecases

import (
	"context"
	"fmt"

	"pr-manager-service/internal/domain"
)

func (u *Usecases) GetPullRequestsByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	// NOTE: проверка существования пользователя
	if _, err := u.storage.GetUserShort(ctx, userID); err != nil {
		return nil, fmt.Errorf("GetUserShort: %w", err)
	}

	return u.storage.GetPullRequestsByReviewer(ctx, userID)
}

func (u *Usecases) UpdateUserStatus(
	ctx context.Context,
	userID string,
	isActive bool,
) (domain.User, domain.Team, error) {
	if err := u.storage.UnitOfWork(ctx, func(s Storage) error {
		user, err := s.GetUserShort(ctx, userID)
		if err != nil {
			return fmt.Errorf("GetUserShort: %w", err)
		}

		if user.IsActive == isActive {
			return nil
		}

		if err := u.storage.UpdateUserStatus(ctx, userID, isActive); err != nil {
			return fmt.Errorf("storage.UpdateUserStatus: %w", err)
		}

		if err := u.storage.UserStatusChangesIncrementBatch(ctx, []string{userID}); err != nil {
			return fmt.Errorf("UserStatusChangeIncrementMany: %w", err)
		}

		return nil
	}); err != nil {
		return domain.User{}, domain.Team{}, fmt.Errorf("UnitOfWork: %w", err)
	}

	return u.storage.GetUserFull(ctx, userID)
}
