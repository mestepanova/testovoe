package http_server

import (
	"net/http"
	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/gin-gonic/gin"
)

// Получить PR'ы, где пользователь назначен ревьювером
// (GET /users/getReview)
func (h *HttpServer) GetUsersGetReview(c *gin.Context, params api.GetUsersGetReviewParams) {
	if err := h.validator.Var(params.UserId, idValidationRules); err != nil {
		handleValidationError(c, err, WithUserID(params.UserId))
		return
	}

	prs, err := h.usecases.GetPullRequestsByReviewer(c.Request.Context(), params.UserId)
	if err != nil {
		handleUsecaseError(c, err, WithUserID(params.UserId))
		return
	}

	response := make([]api.PullRequestShort, 0, len(prs))

	for _, pr := range prs {
		response = append(response, api.PullRequestShort{
			AuthorId:        pr.AuthorUserID,
			PullRequestId:   pr.ID,
			PullRequestName: pr.Name,
			Status:          domain.ConvertPullRequestStatusToApi(pr.Status),
		})
	}

	c.JSON(http.StatusOK, api.UsersGetReviewResponse{
		PullRequests: response,
		UserId:       params.UserId,
	})
}

// Установить флаг активности пользователя
// (POST /users/setIsActive)
func (h *HttpServer) PostUsersSetIsActive(c *gin.Context) {
	request := api.PostUsersSetIsActiveJSONBody{}
	if err := c.ShouldBindJSON(&request); err != nil {
		handleParsingError(c, err)
		return
	}

	if err := h.validator.Var(request.UserId, idValidationRules); err != nil {
		handleValidationError(c, err, WithUserID(request.UserId))
		return
	}

	ctx := c.Request.Context()

	user, team, err := h.usecases.UpdateUserStatus(
		ctx,
		request.UserId,
		request.IsActive,
	)
	if err != nil {
		handleUsecaseError(c, err, WithRequest(request))
		return
	}

	response := api.User{
		UserId:   user.ID,
		Username: user.Name,
		TeamName: team.Name,
		IsActive: user.IsActive,
	}

	c.JSON(http.StatusOK, api.SetIsActiveResponse{
		User: response,
	})
}
