package v1

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/domain/models"
)

func (r *Router) initRecommendationsRoutes() {
	recommendations := r.api.Group("/recommendations", r.identifyUser)
	recommendations.GET("", r.getRecommendations)
	recommendations.POST("/send", r.sendRecommendations)
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

func (r *Router) sendRecommendations(c echo.Context) error {
	rand.Seed(time.Now().Unix())
	switch rand.Intn(100) % 3 {
	case 1:
		return c.JSON(http.StatusInternalServerError, newError(c, errors.New("error stub"), errorInternal))
	case 2:
		return c.JSON(http.StatusUnprocessableEntity, newError(c, errors.New("error stub"), errorNotEnoughInformation))
	default:
		return c.JSON(http.StatusOK, newResult(resultEmailSent))
	}
}
