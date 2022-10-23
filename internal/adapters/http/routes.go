package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) initRoutes() {
	s.server.GET("health", s.health)

	// /api/v1/auth/*
	s.initAuthRoutes()

	// /api/v1/profile/*
	s.initProfileRoutes()

	s.initDiseaseRoutes()
}

func (s *Server) health(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}
