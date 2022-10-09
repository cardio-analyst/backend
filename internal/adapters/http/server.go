package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cardio-analyst/backend/internal/ports/service"
)

type Server struct {
	server      *echo.Echo
	authService service.AuthService
	userService service.UserService
}

func NewServer(authService service.AuthService, userService service.UserService) *Server {
	srv := new(Server)

	e := echo.New()

	// hide echo startup banner
	e.HideBanner = true

	e.Use(middleware.RequestID())
	e.Use(RequestsBodiesLogger())
	e.Use(RequestsLogger())
	e.Use(middleware.Recover())

	srv.server = e

	srv.authService = authService
	srv.userService = userService

	srv.initRoutes()

	return srv
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
