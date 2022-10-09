package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
)

func (s *Server) initProfileRoutes() {
	profile := s.server.Group("/api/v1/profile", s.identifyUser)
	profile.GET("/info", s.getProfileInfo)
	profile.POST("/edit", s.editProfileInfo)
}

func (s *Server) getProfileInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	criteria := models.UserCriteria{
		ID: &userID,
	}

	user, err := s.userService.Get(criteria)
	if err != nil {
		return c.JSON(NewInternalErrorResponse(err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) editProfileInfo(c echo.Context) error {
	var reqData models.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(NewParseRequestDataErrorResponse(err))
	}

	reqData.ID = c.Get(ctxKeyUserID).(uint64)

	if err := s.userService.Update(reqData); err != nil {
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

	return c.JSON(NewOKResponse())
}
