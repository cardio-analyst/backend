package service

import "github.com/cardio-analyst/backend/internal/domain/models"

// AnalysisService TODO
type AnalysisService interface {
	// Create TODO
	Create(analysisData models.Analysis) (err error)
	// Update TODO
	Update(analysisData models.Analysis) (err error)
	// FindAll TODO
	FindAll(userID uint64) (analysis []*models.Analysis, err error)
}
