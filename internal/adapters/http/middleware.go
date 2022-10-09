package http

import (
	"errors"
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

// identifyUser TODO
func (s *Server) identifyUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(headerAuthorization)
		if header == "" {
			return c.JSON(NewUnauthorizedResponse(ErrEmptyAuthHeader))
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.JSON(NewUnauthorizedResponse(ErrInvalidAuthHeader))
		}

		token := headerParts[1]
		if token == "" {
			return c.JSON(NewUnauthorizedResponse(ErrTokenIsEmpty))
		}

		userID, err := s.authService.ValidateToken(token)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrWrongToken):
				return c.JSON(NewForbiddenResponse(err))
			case errors.Is(err, serviceErrors.ErrTokenIsExpired):
				return c.JSON(NewUnauthorizedResponse(err))
			default:
				return c.JSON(NewInternalErrorResponse(err))
			}
		}

		c.Set(ctxKeyUserID, userID)

		return next(c)
	}
}
