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

func (s *Storage) GetTeamByName(ctx context.Context, teamName string) (domain.Team, error) {
	q, args, err := s.builder.Select("id", "name").
		From("teams").
		Where(squirrel.Eq{"name": teamName}).
		ToSql()
	if err != nil {
		return domain.Team{}, fmt.Errorf("query builder: %w", err)
	}

	team := domain.Team{}
	if err := s.querier.QueryRow(ctx, q, args...).Scan(&team.ID, &team.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, domain.ErrTeamNotFound
		}

		return domain.Team{}, fmt.Errorf("conn.QueryRow: %w", err)
	}

	return team, nil
}

func (s *Storage) CreateTeam(ctx context.Context, request domain.CreateTeamRequest, teamID string) error {
	insertTeamsQuery, insertTeamsArgs, err := s.builder.Insert("teams").
		Columns("id", "name").
		Values(teamID, request.Name).
		ToSql()

	if err != nil {
		return fmt.Errorf("insert teams query builder: %w", err)
	}

	if _, err := s.querier.Exec(ctx, insertTeamsQuery, insertTeamsArgs...); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrTeamExists
		}

		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetTeamFullByName(ctx context.Context, teamName string) (domain.Team, []domain.User, error) {
	query, args, err := s.builder.Select(
		"t.id as team_id",
		"t.name as team_name",
		"u.id as user_id",
		"u.name as username",
		"u.is_active as is_active",
	).From("teams t").
		LeftJoin("users u on t.id = u.team_id").
		Where(squirrel.Eq{"t.name": teamName}).
		ToSql()

	if err != nil {
		return domain.Team{}, []domain.User{}, fmt.Errorf("query builder: %w", err)
	}

	rows, err := s.querier.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, []domain.User{}, domain.ErrTeamNotFound
		}

		return domain.Team{}, []domain.User{}, fmt.Errorf("conn.Query: %w", err)
	}

	defer rows.Close()

	var (
		zeroValueTeam domain.Team
		team          domain.Team
		users         []domain.User
	)

	for rows.Next() {
		var (
			teamID   string
			teamName string

			userID   sql.NullString
			username sql.NullString
			isActive sql.NullBool
		)

		if err := rows.Scan(
			&teamID,
			&teamName,
			&userID,
			&username,
			&isActive,
		); err != nil {
			return domain.Team{}, []domain.User{}, fmt.Errorf("rows.Scan: %w", err)
		}

		if team == zeroValueTeam {
			team = domain.Team{
				ID:   teamID,
				Name: teamName,
			}
		}

		if userID.Valid {
			users = append(users, domain.User{
				ID:       userID.String,
				Name:     username.String,
				IsActive: isActive.Bool,
				TeamID:   teamID,
			})
		}
	}

	if team == zeroValueTeam {
		return domain.Team{}, []domain.User{}, domain.ErrTeamNotFound
	}

	return team, users, nil
}
