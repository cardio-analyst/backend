package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// AnalysisService TODO
type AnalysisService interface {
	// Create TODO
	Create(analysisData models.Analysis) (err error)
	// Update TODO
	Update(analysisData models.Analysis) (err error)
	// FindAll TODO
	FindAll(userID uint64) (analysisDataList []*models.Analysis, err error)
}
