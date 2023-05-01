package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
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
	// /:userRole/auth/*
	r.initAuthRoutes()

	// /:userRole/profile/*
	r.initProfileRoutes()

	// /:userRole/feedback/*
	r.initFeedbackRoutes()

	// /customer/*
	r.initCustomerRoutes()

	// /administrator/*
	r.initAdministratorRoutes()
}

func (r *Router) initCustomerRoutes() {
	customerAPI := r.api.Group("/customer")
	{
		// /diseases/*
		r.initDiseasesRoutes(customerAPI)

		// /analyses/*
		r.initAnalysesRoutes(customerAPI)

		// /lifestyles/*
		r.initLifestylesRoutes(customerAPI)

		// /basicIndicators/*
		r.initBasicIndicatorsRoutes(customerAPI)

		// /score/*
		r.initScoreRoutes(customerAPI)

		// /recommendations/*
		r.initRecommendationsRoutes(customerAPI)

		// /tests/*
		r.initQuestionnaireRoutes(customerAPI)
	}
}

func (r *Router) initAdministratorRoutes() {
	adminAPI := r.api.Group("/administrator")
	{
		adminAPI.POST("/secretKey", r.generateSecretKey, r.identifyUser, r.verifyAdministrator)
	}
}
