package http_server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func WithTeamName(name string) slog.Attr {
	return slog.String("team_name", name)
}

func WithRequest(request any) slog.Attr {
	return slog.Any("request", request)
}

func WithUserID(userID string) slog.Attr {
	return slog.String("user_id", userID)
}

func joinAttrs(err error, logAttrs ...slog.Attr) []any {
	attrs := make([]any, 0, len(logAttrs)+1)

	for _, attr := range logAttrs {
		attrs = append(attrs, attr)
	}

	return append(attrs, slog.Any("error", err))
}

func handleValidationError(c *gin.Context, err error, logAttrs ...slog.Attr) {
	attrs := joinAttrs(err, logAttrs...)
	slog.Info("validation error", attrs...)

	message := "invalid request"

	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		parts := make([]string, 0, len(vErrs))
		for _, fe := range vErrs {
			parts = append(parts, validationErrorToText(fe))
		}
		message = strings.Join(parts, "; ")
	}

	c.JSON(http.StatusBadRequest, errorResponse(api.VALIDATIONERR, message))
}

func validationErrorToText(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("field %s is required", field)

	case "min":
		kind := fe.Kind()
		switch kind {
		case reflect.Slice, reflect.Array, reflect.Map:
			return fmt.Sprintf("field %s must contain at least %s item(s)", field, fe.Param())
		default:
			return fmt.Sprintf("field %s must be at least %s", field, fe.Param())
		}

	case "max":
		return fmt.Sprintf("field %s must be at most %s", field, fe.Param())

	default:
		return fmt.Sprintf("field %s is invalid", field)
	}
}

func handleParsingError(c *gin.Context, err error, logAttrs ...slog.Attr) {
	attrs := joinAttrs(err, logAttrs...)
	slog.Info("request parsing error", attrs...)

	msg := "invalid JSON in request body"

	c.JSON(http.StatusBadRequest, errorResponse(api.VALIDATIONERR, msg))
}

func handleUsecaseError(c *gin.Context, err error, logAttrs ...slog.Attr) {
	attrs := joinAttrs(err, logAttrs...)

	var (
		logMessage string
		httpCode   int
		errorResp  api.ErrorResponse

		errNotInTeam domain.ErrNotInTeam
	)

	switch {
	case errors.Is(err, domain.ErrTeamExists):
		logMessage = "team exists error"
		httpCode = http.StatusBadRequest
		errorResp = errorResponse(api.TEAMEXISTS, domain.ErrTeamExists.Error())

	case errors.Is(err, domain.ErrTeamNotFound):
		logMessage = "team not found"
		httpCode = http.StatusNotFound
		errorResp = errorResponse(api.NOTFOUND, domain.ErrTeamNotFound.Error())

	case errors.Is(err, domain.ErrUserNotFound):
		logMessage = "user not found"
		httpCode = http.StatusNotFound
		errorResp = errorResponse(api.NOTFOUND, domain.ErrUserNotFound.Error())

	case errors.Is(err, domain.ErrPullRequestNotFound):
		logMessage = "PR not found"
		httpCode = http.StatusNotFound
		errorResp = errorResponse(api.NOTFOUND, domain.ErrPullRequestNotFound.Error())

	case errors.Is(err, domain.ErrPRExists):
		logMessage = "PR exists"
		httpCode = http.StatusConflict
		errorResp = errorResponse(api.PREXISTS, domain.ErrPRExists.Error())

	case errors.Is(err, domain.ErrPRMerged):
		logMessage = "PR is merged"
		httpCode = http.StatusConflict
		errorResp = errorResponse(api.PRMERGED, domain.ErrPRMerged.Error())

	case errors.Is(err, domain.ErrNotAssigned):
		logMessage = "reviewer not assigned"
		httpCode = http.StatusConflict
		errorResp = errorResponse(api.NOTASSIGNED, domain.ErrNotAssigned.Error())

	case errors.Is(err, domain.ErrNoCandidate):
		logMessage = "no replacement candidate"
		httpCode = http.StatusConflict
		errorResp = errorResponse(api.NOCANDIDATE, domain.ErrNoCandidate.Error())

	case errors.As(err, &errNotInTeam):
		logMessage = "users not in team"
		httpCode = 400
		errorResp = errorResponse(api.NOTINTEAM, errNotInTeam.Error())
		attrs = append(
			attrs,
			slog.Any("user_ids_not_in_team", errNotInTeam.UserIDs),
			slog.Any("team_name", errNotInTeam.TeamName),
		)

	case errors.Is(err, domain.ErrUserInactive):
		logMessage = "inactive user cannot create a pull request"
		httpCode = http.StatusForbidden
		errorResp = errorResponse(api.NOCANDIDATE, domain.ErrUserInactive.Error())

	default:
		slog.Error(err.Error(), attrs...)
		c.JSON(http.StatusInternalServerError, errorResponse(api.INTERNALERR, "internal server error"))
		return
	}

	slog.Info(logMessage, attrs...)
	c.JSON(httpCode, errorResp)
}

func errorResponse(code api.ErrorCode, msg string) api.ErrorResponse {
	return api.ErrorResponse{
		Error: api.Error{
			Code:    code,
			Message: msg,
		},
	}
}
