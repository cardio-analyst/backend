package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) initProfileRoutes() {
	profile := s.server.Group("/api/v1/profile")
	profile.GET("/info", s.getProfileInfo)
	profile.POST("/edit", s.editProfileInfo)
}

func (s *Server) getProfileInfo(c echo.Context) error {
	return c.String(http.StatusOK, "getProfileInfo")
}

func (s *Server) editProfileInfo(c echo.Context) error {
	return c.String(http.StatusOK, "editProfileInfo")
}
