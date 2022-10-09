package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) initProfileRoutes() {
	profile := s.server.Group("/api/v1/profile", s.identifyUser)
	profile.GET("/info", s.getProfileInfo)
	profile.POST("/edit", s.editProfileInfo)
}

func (s *Server) getProfileInfo(c echo.Context) error {
	userLogin := c.Get(ctxKeyUserLogin).(string)
	return c.String(http.StatusOK, fmt.Sprintf("getProfileInfo: %v", userLogin))
}

func (s *Server) editProfileInfo(c echo.Context) error {
	userLogin := c.Get(ctxKeyUserLogin).(string)
	return c.String(http.StatusOK, fmt.Sprintf("editProfileInfo: %v", userLogin))
}
