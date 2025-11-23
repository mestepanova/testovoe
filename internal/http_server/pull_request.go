package http_server

import (
	"errors"
	"net/http"

	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/gin-gonic/gin"
)

const (
	idValidationRules   = "required,min=1,max=36"
	nameValidationRules = "required,min=2,max=50"
)

// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// (POST /pullRequest/create)
func (h *HttpServer) PostPullRequestCreate(c *gin.Context) {
	apiRequest := api.PostPullRequestCreateJSONBody{}
	if err := c.ShouldBindJSON(&apiRequest); err != nil {
		handleParsingError(c, err)
		return
	}

	domainRequest := domain.CreatePullRequestRequest{
		AuthorUserID: apiRequest.AuthorId,
		Name:         apiRequest.PullRequestName,
		ID:           apiRequest.PullRequestId,
	}

	if err := h.validator.Struct(domainRequest); err != nil {
		handleValidationError(c, err, WithRequest(apiRequest))
		return
	}

	pullRequestDomain, err := h.usecases.CreatePullRequest(c.Request.Context(), domainRequest)
	if err != nil {
		handleUsecaseError(c, err, WithRequest(apiRequest))
		return
	}

	c.JSON(http.StatusCreated, api.CreatePullRequestResponse{
		Pr: domain.ConvertPullRequest(pullRequestDomain),
	})
}

// Пометить PR как MERGED (идемпотентная операция)
// (POST /pullRequest/merge)
func (h *HttpServer) PostPullRequestMerge(c *gin.Context) {
	request := api.PostPullRequestMergeJSONBody{}
	if err := c.ShouldBindJSON(&request); err != nil {
		handleParsingError(c, err)
		return
	}

	if err := h.validator.Var(request.PullRequestId, idValidationRules); err != nil {
		handleValidationError(c, err, WithRequest(request))
		return
	}

	pullRequestDomain, err := h.usecases.MergePullRequest(c.Request.Context(), request.PullRequestId)
	if err != nil {
		handleUsecaseError(c, err, WithRequest(request))
		return
	}

	c.JSON(http.StatusOK, api.MergePullRequestResponse{
		Pr: domain.ConvertPullRequest(pullRequestDomain),
	})
}

// Переназначить конкретного ревьювера на другого из его команды
// (POST /pullRequest/reassign)
func (h *HttpServer) PostPullRequestReassign(c *gin.Context) {
	request := api.PostPullRequestReassignJSONBody{}
	if err := c.ShouldBindJSON(&request); err != nil {
		handleParsingError(c, err)
		return
	}

	if err := errors.Join(
		h.validator.Var(request.OldUserId, idValidationRules),
		h.validator.Var(request.PullRequestId, idValidationRules),
	); err != nil {
		handleValidationError(c, err, WithRequest(request))
		return
	}

	pullRequestDomain, newReviewerID, err := h.usecases.ReassignPullRequest(
		c.Request.Context(),
		request.PullRequestId,
		request.OldUserId,
	)
	if err != nil {
		handleUsecaseError(c, err, WithRequest(request))
		return
	}

	c.JSON(http.StatusOK, api.ReassignPullRequestResponse{
		Pr:         domain.ConvertPullRequest(pullRequestDomain),
		ReplacedBy: newReviewerID,
	})
}
