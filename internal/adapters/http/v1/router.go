package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/ports/service"
)

type Router struct {
	api      *echo.Group
	services service.Services
}

func NewRouter(api *echo.Group, services service.Services) *Router {
	return &Router{
		api:      api,
		services: services,
	}
}

func (r *Router) InitRoutes() {
	// /auth/*
	r.initAuthRoutes()

	// /profile/*
	r.initProfileRoutes()

	// /diseases/*
	r.initDiseasesRoutes()

	// /analyses/*
	r.initAnalysesRoutes()
}
