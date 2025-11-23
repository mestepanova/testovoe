package domain

import (
	"errors"
)

var (
	ErrTeamNotFound        = errors.New("team not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrPullRequestNotFound = errors.New("pull request not found")
	ErrTeamExists          = errors.New("team_name already exists")
	ErrPRExists            = errors.New("PR id already exists")
	ErrPRMerged            = errors.New("cannot reassign on merged PR")
	ErrNotAssigned         = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrUserInactive        = errors.New("inactive user cannot create a pull request")
	ErrInternal            = errors.New("internal server error")
)

type ErrNotInTeam struct {
	UserIDs  []string
	TeamName string
}

func NewErrNotInTeam(teamName string, userIDs []string) ErrNotInTeam {
	return ErrNotInTeam{
		UserIDs:  userIDs,
		TeamName: teamName,
	}
}

func (e ErrNotInTeam) Error() string {
	return "some users are not in team"
}
