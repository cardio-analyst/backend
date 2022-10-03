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
	auth.POST("/sign-up", s.signUp)
	auth.POST("/sign-in", s.signIn)
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

type signInRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

// signIn TODO
func (s *Server) signIn(c echo.Context) error {
	var req signInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(NewParseRequestDataErrorResponse(err))
	}

	credentials := models.UserCredentials{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := s.authService.GetToken(credentials)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidUserCredentials):
			return c.JSON(NewInvalidRequestDataResponse(err))
		case errors.Is(err, serviceErrors.ErrWrongCredentials):
			return c.JSON(NewForbiddenResponse(err))
		default:
			return c.JSON(NewInternalErrorResponse(err))
		}
	}

	res := &signInResponse{
		Token: token,
	}

	return c.JSON(http.StatusOK, res)
}
