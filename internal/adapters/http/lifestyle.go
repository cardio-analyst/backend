package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/domain/models"
)

func (s *Server) initLifestylesRoutes() {
	lifestyle := s.server.Group("/api/v1/lifestyle", s.identifyUser)
	lifestyle.GET("/info", s.getLifestyleInfo)
	lifestyle.PUT("/edit", s.editLifestyleInfo)
}

func (s *Server) getLifestyleInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	userLifestyle, err := s.services.Lifestyle().Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, userLifestyle)
}

func (s *Server) editLifestyleInfo(c echo.Context) error {
	var reqData models.Lifestyle
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := s.services.Lifestyle().Update(reqData); err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, NewResult(resultUpdated))
}
