package http

import (
	"errors"
	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) initDiseaseRoutes() {
	disease := s.server.Group("/api/v1/disease", s.identifyUser)
	disease.GET("/info", s.getDiseaseInfo)
	disease.PUT("/edit", s.editDiseaseInfo)
}

func (s *Server) getDiseaseInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	disease, err := s.services.Disease().Get(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, disease)
}

func (s *Server) editDiseaseInfo(c echo.Context) error {
	var reqData models.Disease
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := s.services.Disease().Update(reqData); err != nil {
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

	return c.JSON(http.StatusOK, NewResult(resultUpdated))
}
