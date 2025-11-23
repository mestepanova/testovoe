package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pr-manager-service/internal/domain"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func (s *Storage) UpdateUserStatus(ctx context.Context, userID string, isActive bool) error {
	query, args, err := s.builder.Update("users").
		Set("is_active", isActive).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}

	_, err = s.querier.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetUserFull(ctx context.Context, userID string) (domain.User, domain.Team, error) {
	query, args, err := s.builder.Select(
		"t.id as team_id",
		"t.name as team_name",
		"u.name as username",
		"u.is_active as is_active",
	).From("users u").
		LeftJoin("teams t on t.id = u.team_id").
		Where(squirrel.Eq{"u.id": userID}).
		ToSql()

	if err != nil {
		return domain.User{}, domain.Team{}, fmt.Errorf("query builder: %w", err)
	}

	var (
		teamID   sql.NullString
		teamName sql.NullString
		username string
		isActive bool
	)

	if err := s.querier.QueryRow(ctx, query, args...).Scan(
		&teamID,
		&teamName,
		&username,
		&isActive,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.Team{}, domain.ErrUserNotFound
		}

		return domain.User{}, domain.Team{}, fmt.Errorf("conn.QueryRow: %w", err)
	}

	if !teamID.Valid {
		return domain.User{}, domain.Team{}, domain.ErrTeamNotFound
	}

	team := domain.Team{
		ID:   teamID.String,
		Name: teamName.String,
	}

	user := domain.User{
		ID:       userID,
		Name:     username,
		TeamID:   teamID.String,
		IsActive: isActive,
	}

	return user, team, nil
}

func (s *Storage) GetUserShort(ctx context.Context, userID string) (domain.User, error) {
	query, args, err := s.builder.Select("id, name, is_active, team_id").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return domain.User{}, fmt.Errorf("query builder: %w", err)
	}

	user := domain.User{}
	if err := s.querier.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.IsActive,
		&user.TeamID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("conn.QueryRow: %w", err)
	}

	return user, nil
}

func (s *Storage) GetActiveColleagues(ctx context.Context, userID string) ([]domain.User, error) {
	teamSubqueryQuery, _, err := s.builder.Select("team_id").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("subquery query builder: %w", err)
	}

	query, args, err := s.builder.Select(
		"u.id as user_id",
		"u.name as username",
		"u.is_active as is_active",
		"u.team_id as team_id",
	).From("users u").
		Where(squirrel.And{
			squirrel.Expr("u.team_id = (" + teamSubqueryQuery + ")"),
			squirrel.NotEq{"u.id": userID},
			squirrel.Eq{"u.is_active": true},
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("select query builder: %w", err)
	}

	rows, err := s.querier.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("conn.Query: %w", err)
	}

	users := []domain.User{}
	for rows.Next() {
		var user domain.User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.IsActive,
			&user.TeamID,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) CreateUsers(ctx context.Context, requests []domain.CreateUserRequest, teamID string) error {
	builder := s.builder.Insert("users").
		Columns("id", "name", "is_active", "team_id")

	for _, member := range requests {
		builder = builder.Values(
			member.ID,
			member.Name,
			member.IsActive,
			teamID,
		)
	}

	builder = builder.Suffix(`
		on conflict (id) do update set 
        name = excluded.name,
        is_active = excluded.is_active,
		team_id = excluded.team_id
	`)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("insert users query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}
