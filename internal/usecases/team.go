package usecases

import (
	"context"
	"fmt"

	"pr-manager-service/internal/domain"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (u *Usecases) CreateTeam(ctx context.Context, request domain.CreateTeamRequest) error {
	if err := u.storage.UnitOfWork(ctx, func(s Storage) error {
		teamID := uuid.NewString()

		if err := u.storage.CreateTeam(ctx, request, teamID); err != nil {
			return fmt.Errorf("CreateTeam: %w", err)
		}

		if err := u.storage.CreateUsers(ctx, request.Members, teamID); err != nil {
			return fmt.Errorf("CreateUsers: %w", err)
		}

		userIDs := lo.Map(request.Members, func(user domain.CreateUserRequest, _ int) string {
			return user.ID
		})
		if err := u.storage.UserStatsCreateBatch(ctx, userIDs); err != nil {
			return fmt.Errorf("UserStatsCreateBatch: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("UnitOfWork: %w", err)
	}

	return nil
}

func (u *Usecases) GetTeamFullByName(ctx context.Context, teamName string) (domain.Team, []domain.User, error) {
	return u.storage.GetTeamFullByName(ctx, teamName)
}
