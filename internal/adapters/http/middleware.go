package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
)

const (
	headerAuthorization = "Authorization"
	ctxKeyUserID        = "userID"
)

var (
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrTokenIsEmpty      = errors.New("token is empty")
)

func (s *Server) identifyUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(headerAuthorization)
		if header == "" {
			return c.JSON(http.StatusUnauthorized, NewError(c, ErrEmptyAuthHeader, errorWrongAuthHeader))
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, NewError(c, ErrInvalidAuthHeader, errorWrongAuthHeader))
		}

		token := headerParts[1]
		if token == "" {
			return c.JSON(http.StatusUnauthorized, NewError(c, ErrTokenIsEmpty, errorWrongAuthHeader))
		}

		userID, err := s.authService.ValidateAccessToken(token)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrWrongToken):
				return c.JSON(http.StatusBadRequest, NewError(c, err, errorWrongToken))
			case errors.Is(err, serviceErrors.ErrTokenIsExpired):
				return c.JSON(http.StatusUnauthorized, NewError(c, err, errorTokenExpired))
			default:
				return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
			}
		}

		c.Set(ctxKeyUserID, userID)

		return next(c)
	}
}
