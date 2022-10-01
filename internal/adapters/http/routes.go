package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// initRoutes TODO
func (s *Server) initRoutes() {
	s.server.GET("health", s.health)

	// /api/v1/auth/*
	s.initAuthRoutes()
}

// health TODO
func (s *Server) health(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}
