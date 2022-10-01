package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
)

// initAuthRoutes TODO
func (s *Server) initAuthRoutes() {
	auth := s.server.Group("/api/v1/auth")
	auth.POST("/sign-up", s.signUp)
	auth.POST("/sign-in", s.signIn)
	auth.POST("/logout", s.logout)
}

// signUp TODO
func (s *Server) signUp(c echo.Context) error {
	var reqData models.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(NewParseRequestDataErrorResponse(err))
	}

	if err := s.authService.RegisterUser(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidUserData):
			return c.JSON(NewInvalidRequestDataResponse(err))
		case errors.Is(err, serviceErrors.ErrUserLoginAlreadyOccupied):
			return c.JSON(NewAlreadyRegisteredWithLoginResponse(reqData.Login))
		case errors.Is(err, serviceErrors.ErrUserEmailAlreadyOccupied):
			return c.JSON(NewAlreadyRegisteredWithEmailResponse(reqData.Email))
		default:
			return c.JSON(NewInternalErrorResponse(err))
		}
	}

	return c.JSON(http.StatusOK, nil)
}

// signIn TODO
func (s *Server) signIn(c echo.Context) error {
	return c.String(http.StatusOK, "sign-in")
}

// logout TODO
func (s *Server) logout(c echo.Context) error {
	return c.String(http.StatusOK, "logout")
}
