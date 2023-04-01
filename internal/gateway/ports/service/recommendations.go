package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

type RecommendationsService interface {
	GetRecommendations(userID uint64) (recommendations []*models.Recommendation, err error)
}
