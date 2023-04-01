package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// BasicIndicatorsService TODO
type BasicIndicatorsService interface {
	// Create TODO
	Create(basicIndicatorsData models.BasicIndicators) (err error)
	// Update TODO
	Update(basicIndicatorsData models.BasicIndicators) (err error)
	// FindAll TODO
	FindAll(userID uint64) (basicIndicatorsDataList []*models.BasicIndicators, err error)
}
