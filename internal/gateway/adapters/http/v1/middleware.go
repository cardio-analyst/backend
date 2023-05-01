package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const headerAuthorization = "Authorization"

const (
	ctxKeyUserID   = "userID"
	ctxKeyUserRole = "userRole"
)

const (
	userRolePathKey                   = "userRole"
	userRoleCustomerPathParamKey      = "customer"
	userRoleModeratorPathParamKey     = "moderator"
	userRoleAdministratorPathParamKey = "administrator"
)

// possible middleware error designations
const (
	errorWrongAuthHeader    = "WrongAuthHeader"
	errorAccessTokenExpired = "AccessTokenExpired"
	errorWrongAccessToken   = "WrongAccessToken"
	errorForbiddenByRole    = "ForbiddenByRole"
)

var (
	errEmptyAuthHeader   = errors.New("empty auth header")
	errInvalidAuthHeader = errors.New("invalid auth header")
	errTokenIsEmpty      = errors.New("token is empty")
	errForbiddenByRole   = errors.New("forbidden by role")
)

func (r *Router) identifyUser(next echo.HandlerFunc) echo.HandlerFunc {
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

		c.Set(ctxKeyUserID, userID)
		c.Set(ctxKeyUserRole, userRole)

		return next(c)
	}
}

func (r *Router) verifyCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get(ctxKeyUserRole).(model.UserRole)
		if userRole != model.UserRoleCustomer {
			err := fmt.Errorf("%w: %q", errForbiddenByRole, userRole)
			return c.JSON(http.StatusForbidden, newError(c, err, errorForbiddenByRole))
		}
		return next(c)
	}
}

func (r *Router) verifyModerator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get(ctxKeyUserRole).(model.UserRole)
		if userRole != model.UserRoleModerator {
			err := fmt.Errorf("%w: %q", errForbiddenByRole, userRole)
			return c.JSON(http.StatusForbidden, newError(c, err, errorForbiddenByRole))
		}
		return next(c)
	}
}

func (r *Router) verifyAdministrator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get(ctxKeyUserRole).(model.UserRole)
		if userRole != model.UserRoleAdministrator {
			err := fmt.Errorf("%w: %q", errForbiddenByRole, userRole)
			return c.JSON(http.StatusForbidden, newError(c, err, errorForbiddenByRole))
		}
		return next(c)
	}
}

func (r *Router) parseUserRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var role model.UserRole

		roleStr := c.Param(userRolePathKey)
		switch roleStr {
		case userRoleCustomerPathParamKey:
			role = model.UserRoleCustomer
		case userRoleModeratorPathParamKey:
			role = model.UserRoleModerator
		case userRoleAdministratorPathParamKey:
			role = model.UserRoleAdministrator
		default:
			err := fmt.Errorf("undefined user role: %q", roleStr)
			return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
		}

		userRole, ok := c.Get(ctxKeyUserRole).(model.UserRole)
		if ok {
			if userRole != role {
				return c.JSON(http.StatusForbidden, newError(c, model.ErrForbiddenByRole, errorForbiddenByRole))
			}
		} else {
			c.Set(ctxKeyUserRole, role)
		}

		return next(c)
	}
}
