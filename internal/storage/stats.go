package storage

import (
	"context"
	"fmt"
	"pr-manager-service/internal/domain"
	"time"

	"github.com/Masterminds/squirrel"
)

func (s *Storage) UserStatsCreateBatch(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	builder := s.builder.
		Insert("users_stats").
		Columns("user_id").
		Suffix("on conflict (user_id) do nothing")

	for _, id := range userIDs {
		builder = builder.Values(id)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (s *Storage) PullRequestStatsCreate(ctx context.Context, pullRequestID string, assignmentsCount int) error {
	query, args, err := s.builder.
		Insert("pull_requests_stats").
		Columns("pull_request_id", "assignments_count").
		Values(pullRequestID, assignmentsCount).
		Suffix("on conflict (pull_request_id) do nothing").ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (s *Storage) UserAssignmentsIncrementBatch(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	query, args, err := s.builder.
		Update("users_stats").
		Set("assignments_count", squirrel.Expr("assignments_count + 1")).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"user_id": userIDs}). // user_id IN (...)
		ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (s *Storage) UserStatusChangesIncrementBatch(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	query, args, err := s.builder.
		Update("users_stats").
		Set("status_changes_count", squirrel.Expr("status_changes_count + 1")).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"user_id": userIDs}). // IN (...)
		ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (s *Storage) PullRequestAssignmentsIncrement(ctx context.Context, pullRequestID string) error {
	query, args, err := s.builder.Update("pull_requests_stats").
		Set("assignments_count", squirrel.Expr("assignments_count + 1")).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"pull_request_id": pullRequestID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetUsersStats(ctx context.Context) (userStats []domain.UserStats, err error) {
	userStatsQuery, userStatsArgs, err := s.builder.Select(
		"user_id",
		"assignments_count",
		"status_changes_count",
	).
		From("users_stats").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("query builder: %w", err)
	}

	rows, err := s.querier.Query(ctx, userStatsQuery, userStatsArgs...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		stat := domain.UserStats{}

		if err := rows.Scan(
			&stat.UserID,
			&stat.AssignmentsCount,
			&stat.StatusChangesCount,
		); err != nil {
			return nil, fmt.Errorf("Scan: %w", err)
		}

		userStats = append(userStats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Err: %w", err)
	}

	return userStats, nil
}

func (s *Storage) GetPullRequestsStats(ctx context.Context) (pullRequestsStats []domain.PullRequestStats, err error) {
	pullRequestsStatsQuery, pullRequestsStatsArgs, err := s.builder.Select(
		"pull_request_id",
		"assignments_count",
	).
		From("pull_requests_stats").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("query builder: %w", err)
	}

	rows, err := s.querier.Query(ctx, pullRequestsStatsQuery, pullRequestsStatsArgs...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		stat := domain.PullRequestStats{}

		if err := rows.Scan(
			&stat.PullRequestID,
			&stat.AssignmentsCount,
		); err != nil {
			return nil, fmt.Errorf("Scan: %w", err)
		}

		pullRequestsStats = append(pullRequestsStats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Err: %w", err)
	}

	return pullRequestsStats, nil
}
