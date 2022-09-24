package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) initRoutes() {
	s.server.GET("hello", s.hello)
}

func (s *Server) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
