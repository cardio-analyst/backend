package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
)

func (s *Server) initAuthRoutes() {
	auth := s.server.Group("/api/v1/auth")
	auth.POST("/signUp", s.signUp)
	auth.POST("/signIn", s.signIn)
	auth.POST("/refreshTokens", s.refreshTokens)
}

func (s *Server) signUp(c echo.Context) error {
	var reqData models.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(c, err, errorParseRequestData))
	}

	if err := s.services.Auth().RegisterUser(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidUserData):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorInvalidRequestData))
		case errors.Is(err, serviceErrors.ErrUserLoginAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorLoginAlreadyOccupied))
		case errors.Is(err, serviceErrors.ErrUserEmailAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorEmailAlreadyOccupied))
		default:
			return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, NewResult(resultRegistered))
}

func (s *Server) signIn(c echo.Context) error {
	var reqData models.UserCredentials
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(c, err, errorParseRequestData))
	}

	tokens, err := s.services.Auth().GetTokens(reqData, c.RealIP())
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidUserCredentials):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorInvalidRequestData))
		case errors.Is(err, serviceErrors.ErrWrongCredentials):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorWrongCredentials))
		default:
			return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

type refreshTokensRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) refreshTokens(c echo.Context) error {
	var reqData refreshTokensRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(c, err, errorParseRequestData))
	}

	tokens, err := s.services.Auth().RefreshTokens(reqData.RefreshToken, c.RealIP())
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrWrongToken):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorWrongRefreshToken))
		case errors.Is(err, serviceErrors.ErrTokenIsExpired):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorRefreshTokenExpired))
		case errors.Is(err, serviceErrors.ErrSessionNotFound):
			return c.JSON(http.StatusBadRequest, NewError(c, err, errorWrongRefreshToken))
		case errors.Is(err, serviceErrors.ErrIPIsNotInWhitelist):
			return c.JSON(http.StatusForbidden, NewError(c, err, errorIPNotAllowed))
		default:
			return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}
