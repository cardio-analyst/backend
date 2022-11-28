package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/domain/models"
)

func (r *Router) initRecommendationsRoutes() {
	recommendations := r.api.Group("/recommendations", r.identifyUser)
	recommendations.GET("", r.getRecommendations)
}

type getRecommendationsResponse struct {
	Recommendations []*models.Recommendation `json:"recommendations"`
}

func (r *Router) getRecommendations(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	recommendations, err := r.services.Recommendations().GetRecommendations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &getRecommendationsResponse{
		Recommendations: recommendations,
	})
}
