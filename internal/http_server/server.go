package http_server

import (
	"context"
	"reflect"
	"strings"

	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/go-playground/validator/v10"
)

type usecases interface {
	GetPullRequestsByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error)
	UpdateUserStatus(ctx context.Context, userID string, isActive bool) (domain.User, domain.Team, error)

	CreateTeam(ctx context.Context, team domain.CreateTeamRequest) error
	GetTeamFullByName(ctx context.Context, teamName string) (domain.Team, []domain.User, error)

	CreatePullRequest(ctx context.Context, pr domain.CreatePullRequestRequest) (domain.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (domain.PullRequest, error)
	ReassignPullRequest(ctx context.Context, prID, oldUserID string) (
		pr domain.PullRequest,
		newReviewerID string,
		err error,
	)

	GetStats(ctx context.Context) ([]domain.UserStats, []domain.PullRequestStats, error)
}

var _ api.ServerInterface = (*HttpServer)(nil)

type HttpServer struct {
	usecases  usecases
	validator *validator.Validate
}

func NewHttpServer(usecases usecases) *HttpServer {
	return &HttpServer{
		usecases:  usecases,
		validator: NewValidator(),
	}
}

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}
