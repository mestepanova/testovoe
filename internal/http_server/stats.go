package http_server

import (
	"net/http"
	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func (h *HttpServer) GetStatsGet(c *gin.Context) {
	userStats, pullRequestsStats, err := h.usecases.GetStats(c.Request.Context())
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response := api.Stats{
		PullRequestsStats: lo.Map(
			pullRequestsStats,
			func(pullRequestStat domain.PullRequestStats, _ int) api.PullRequestsStats {
				return api.PullRequestsStats{
					PullRequestId:    pullRequestStat.PullRequestID,
					AssignmentsCount: int(pullRequestStat.AssignmentsCount),
				}
			},
		),
		UserStats: lo.Map(userStats, func(userStat domain.UserStats, _ int) api.UserStats {
			return api.UserStats{
				UserId:             userStat.UserID,
				StatusChangesCount: int(userStat.StatusChangesCount),
				AssignmentsCount:   int(userStat.AssignmentsCount),
			}
		}),
	}

	c.JSON(http.StatusOK, response)
}
