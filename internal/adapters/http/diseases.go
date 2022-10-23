package http

import (
	"net/http"

	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/labstack/echo/v4"
)

func (s *Server) initDiseasesRoutes() {
	diseases := s.server.Group("/api/v1/diseases", s.identifyUser)
	diseases.GET("/info", s.getDiseasesInfo)
	diseases.PUT("/edit", s.editDiseasesInfo)
}

func (s *Server) getDiseasesInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	userDiseases, err := s.services.Diseases().Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, userDiseases)
}

func (s *Server) editDiseasesInfo(c echo.Context) error {
	var reqData models.Diseases
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := s.services.Diseases().Update(reqData); err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, NewResult(resultUpdated))
}
