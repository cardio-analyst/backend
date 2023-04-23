package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/pkg/model"
)

const (
	headerAuthorization = "Authorization"
	ctxKeyUserID        = "userID"
	ctxKeyUserRole      = "userRole"
)

// possible middleware error designations
const (
	errorWrongAuthHeader    = "WrongAuthHeader"
	errorAccessTokenExpired = "AccessTokenExpired"
	errorWrongAccessToken   = "WrongAccessToken"
)

var (
	errEmptyAuthHeader   = errors.New("empty auth header")
	errInvalidAuthHeader = errors.New("invalid auth header")
	errTokenIsEmpty      = errors.New("token is empty")
	errForbiddenByRole   = errors.New("forbidden by role")
)

func (r *Router) identifyCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(headerAuthorization)
		if header == "" {
			return c.JSON(http.StatusUnauthorized, newError(c, errEmptyAuthHeader, errorWrongAuthHeader))
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, newError(c, errInvalidAuthHeader, errorWrongAuthHeader))
		}

		token := headerParts[1]
		if token == "" {
			return c.JSON(http.StatusUnauthorized, newError(c, errTokenIsEmpty, errorWrongAuthHeader))
		}

		userID, userRole, err := r.services.Auth().ValidateAccessToken(c.Request().Context(), token)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrWrongToken):
				return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongAccessToken))
			case errors.Is(err, model.ErrTokenIsExpired):
				return c.JSON(http.StatusUnauthorized, newError(c, err, errorAccessTokenExpired))
			default:
				return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
			}
		}

		if userRole != model.UserRoleCustomer {
			err = fmt.Errorf("%w: %q", errForbiddenByRole, userRole)
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongAccessToken))
		}

		c.Set(ctxKeyUserID, userID)
		c.Set(ctxKeyUserRole, userRole)

		return next(c)
	}
}
