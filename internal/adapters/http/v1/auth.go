package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
)

// possible auth errors designations
const (
	errorInvalidFirstName    = "InvalidFirstName"
	errorInvalidLastName     = "InvalidLastName"
	errorInvalidRegion       = "InvalidRegion"
	errorInvalidBirthDate    = "InvalidBirthDate"
	errorInvalidLogin        = "InvalidLogin"
	errorInvalidEmail        = "InvalidEmail"
	errorInvalidPassword     = "InvalidPassword"
	errorRefreshTokenExpired = "RefreshTokenExpired"
	errorWrongRefreshToken   = "WrongRefreshToken"
	errorWrongCredentials    = "WrongCredentials"
	errorIPNotAllowed        = "IPNotAllowed"
)

func (r *Router) initAuthRoutes() {
	auth := r.api.Group("/auth")
	auth.POST("/signUp", r.signUp)
	auth.POST("/signIn", r.signIn)
	auth.POST("/refreshTokens", r.refreshTokens)
}

func (r *Router) signUp(c echo.Context) error {
	var reqData models.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	if err := r.services.Auth().RegisterUser(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidFirstName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidFirstName))
		case errors.Is(err, serviceErrors.ErrInvalidLastName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLastName))
		case errors.Is(err, serviceErrors.ErrInvalidRegion):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRegion))
		case errors.Is(err, serviceErrors.ErrInvalidBirthDate):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidBirthDate))
		case errors.Is(err, serviceErrors.ErrInvalidLogin):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLogin))
		case errors.Is(err, serviceErrors.ErrInvalidEmail):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidEmail))
		case errors.Is(err, serviceErrors.ErrInvalidPassword):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidPassword))
		case errors.Is(err, serviceErrors.ErrInvalidUserData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, serviceErrors.ErrUserLoginAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorLoginAlreadyOccupied))
		case errors.Is(err, serviceErrors.ErrUserEmailAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorEmailAlreadyOccupied))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultRegistered))
}

func (r *Router) signIn(c echo.Context) error {
	var reqData models.UserCredentials
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	tokens, err := r.services.Auth().GetTokens(reqData, c.RealIP())
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidUserCredentials):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, serviceErrors.ErrWrongCredentials):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongCredentials))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

type refreshTokensRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r *Router) refreshTokens(c echo.Context) error {
	var reqData refreshTokensRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	tokens, err := r.services.Auth().RefreshTokens(reqData.RefreshToken, c.RealIP())
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrWrongToken):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongRefreshToken))
		case errors.Is(err, serviceErrors.ErrTokenIsExpired):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorRefreshTokenExpired))
		case errors.Is(err, serviceErrors.ErrSessionNotFound):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongRefreshToken))
		case errors.Is(err, serviceErrors.ErrIPIsNotInWhitelist):
			return c.JSON(http.StatusForbidden, newError(c, err, errorIPNotAllowed))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}
