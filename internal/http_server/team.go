package http_server

import (
	"net/http"
	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/gin-gonic/gin"
)

// Создать команду с участниками (создаёт/обновляет пользователей)
// (POST /team/add)
func (h *HttpServer) PostTeamAdd(c *gin.Context) {
	apiRequest := api.Team{}
	if err := c.ShouldBindJSON(&apiRequest); err != nil {
		handleParsingError(c, err)
		return
	}

	domainRequest := domain.CreateTeamRequest{
		Name:    apiRequest.TeamName,
		Members: make([]domain.CreateUserRequest, 0, len(apiRequest.Members)),
	}

	for _, member := range apiRequest.Members {
		domainRequest.Members = append(domainRequest.Members, domain.CreateUserRequest{
			ID:       member.UserId,
			Name:     member.Username,
			IsActive: member.IsActive,
		})
	}

	if err := h.validator.Struct(domainRequest); err != nil {
		handleValidationError(c, err, WithRequest(apiRequest))
		return
	}

	if err := h.usecases.CreateTeam(c.Request.Context(), domainRequest); err != nil {
		handleUsecaseError(c, err, WithRequest(apiRequest))
		return
	}

	c.JSON(http.StatusCreated, api.TeamAddResponse{
		Team: apiRequest,
	})
}

// Получить команду с участниками
// (GET /team/get)
func (h *HttpServer) GetTeamGet(c *gin.Context, params api.GetTeamGetParams) {
	if err := h.validator.Var(params.TeamName, nameValidationRules); err != nil {
		handleValidationError(c, err, WithTeamName(params.TeamName))
		return
	}

	team, users, err := h.usecases.GetTeamFullByName(c.Request.Context(), params.TeamName)
	if err != nil {
		handleUsecaseError(c, err, WithTeamName(params.TeamName))
		return
	}

	response := api.Team{
		TeamName: team.Name,
		Members:  make([]api.TeamMember, 0, len(users)),
	}

	for _, user := range users {
		response.Members = append(response.Members, api.TeamMember{
			IsActive: user.IsActive,
			UserId:   user.ID,
			Username: user.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}
