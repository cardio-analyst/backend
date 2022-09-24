package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	server      *echo.Echo
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	srv := new(Server)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	srv.server = e
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
