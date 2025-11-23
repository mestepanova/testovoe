package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pr-manager-service/internal/domain"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func (s *Storage) UpdatePullRequestReviewersIDs(ctx context.Context, prID string, reviewersIDs []string) error {
	query, args, err := s.builder.Update("pull_requests").
		Set("reviewers_ids", reviewersIDs).
		Where(squirrel.Eq{"id": prID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	if _, err = s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetPullRequestsByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	query, args, err := s.builder.Select(
		"id",
		"author_id",
		"reviewers_ids",
		"name",
		"created_at",
		"merged_at",
		"status",
	).From("pull_requests").
		Where(squirrel.Expr("? = ANY(reviewers_ids)", userID)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("query builder: %w", err)
	}

	rows, err := s.querier.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("conn.Query: %w", err)
	}

	pullRequests := make([]domain.PullRequest, 0, 10)
	for rows.Next() {
		pullRequest := domain.PullRequest{}

		if err = rows.Scan(
			&pullRequest.ID,
			&pullRequest.AuthorUserID,
			&pullRequest.ReviewersUsersIDs,
			&pullRequest.Name,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt,
			&pullRequest.Status,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil
}

func (s *Storage) GetPullRequestByID(ctx context.Context, prID string) (domain.PullRequest, error) {
	query, args, err := s.builder.Select(
		"id",
		"author_id",
		"reviewers_ids",
		"name",
		"created_at",
		"merged_at",
		"status",
	).From("pull_requests").
		Where(squirrel.Eq{"id": prID}).
		ToSql()

	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("query builder: %w", err)
	}

	pullRequest := domain.PullRequest{}

	if err := s.querier.QueryRow(ctx, query, args...).Scan(
		&pullRequest.ID,
		&pullRequest.AuthorUserID,
		&pullRequest.ReviewersUsersIDs,
		&pullRequest.Name,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
		&pullRequest.Status,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PullRequest{}, domain.ErrPullRequestNotFound
		}

		return domain.PullRequest{}, fmt.Errorf("conn.QueryRow: %w", err)
	}

	return pullRequest, nil
}

func (s *Storage) CreatePullRequest(
	ctx context.Context,
	request domain.CreatePullRequestRequest,
	reviewersIDs []string,
) (pr domain.PullRequest, err error) {
	timeNow := time.Now()
	pr.ID = request.ID
	pr.AuthorUserID = request.AuthorUserID
	pr.ReviewersUsersIDs = reviewersIDs
	pr.Name = request.Name
	pr.CreatedAt = &timeNow
	pr.Status = domain.StatusOpen

	query, args, err := s.builder.Insert("pull_requests").
		Columns(
			"id",
			"author_id",
			"reviewers_ids",
			"name",
			"created_at",
			"merged_at",
			"status",
		).
		Values(
			pr.ID,
			pr.AuthorUserID,
			pr.ReviewersUsersIDs,
			pr.Name,
			pr.CreatedAt,
			pr.MergedAt,
			pr.Status,
		).
		ToSql()

	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		if isUniqueViolation(err) {
			return domain.PullRequest{}, domain.ErrPRExists
		}

		return domain.PullRequest{}, fmt.Errorf("tx.Exec: %w", err)
	}

	return pr, nil
}

func (s *Storage) UpdatePullRequestStatus(ctx context.Context, prID string, newStatus domain.PullRequestStatus) error {
	builder := s.builder.Update("pull_requests").
		Set("status", newStatus).
		Where(squirrel.Eq{"id": prID})

	if newStatus == domain.StatusMerged {
		builder = builder.Set("merged_at", time.Now())
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	_, err = s.querier.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}
