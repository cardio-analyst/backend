package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cardio-analyst/backend/internal/gateway/adapters/http/v1"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
)

type Server struct {
	server *echo.Echo
}

func NewServer(services service.Services) *Server {
	e := echo.New()

	// hide echo startup banner
	e.HideBanner = true

	e.Use(middleware.RequestID())
	e.Use(RequestsBodiesLogger())
	e.Use(RequestsLogger())
	e.Use(middleware.Recover())

	e.GET("health", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	api := e.Group("/api/v1")
	r := v1.NewRouter(api, services)
	r.InitRoutes()

	return &Server{
		server: e,
	}
}

func (s *Server) Start(listenAddress string) error {
	err := s.server.Start(listenAddress)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	shutdownCtx := context.Background()
	return s.server.Shutdown(shutdownCtx)
}
