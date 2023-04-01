package v1

import (
	"errors"
	errors2 "github.com/cardio-analyst/backend/internal/gateway/domain/errors"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"net/http"

	"github.com/labstack/echo/v4"
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
		case errors.Is(err, errors2.ErrInvalidFirstName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidFirstName))
		case errors.Is(err, errors2.ErrInvalidLastName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLastName))
		case errors.Is(err, errors2.ErrInvalidRegion):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRegion))
		case errors.Is(err, errors2.ErrInvalidBirthDate):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidBirthDate))
		case errors.Is(err, errors2.ErrInvalidLogin):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLogin))
		case errors.Is(err, errors2.ErrInvalidEmail):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidEmail))
		case errors.Is(err, errors2.ErrInvalidPassword):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidPassword))
		case errors.Is(err, errors2.ErrInvalidUserData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, errors2.ErrUserLoginAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorLoginAlreadyOccupied))
		case errors.Is(err, errors2.ErrUserEmailAlreadyOccupied):
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
		case errors.Is(err, errors2.ErrInvalidUserCredentials):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, errors2.ErrWrongCredentials):
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
		case errors.Is(err, errors2.ErrWrongToken):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongRefreshToken))
		case errors.Is(err, errors2.ErrTokenIsExpired):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorRefreshTokenExpired))
		case errors.Is(err, errors2.ErrSessionNotFound):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongRefreshToken))
		case errors.Is(err, errors2.ErrIPIsNotInWhitelist):
			return c.JSON(http.StatusForbidden, newError(c, err, errorIPNotAllowed))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}
