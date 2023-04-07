package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type RecommendationsService interface {
	GetRecommendations(userID uint64) (recommendations []*domain.Recommendation, err error)
}
