package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// LifestyleService TODO
type LifestyleService interface {
	// Update TODO
	Update(lifestyleData models.Lifestyle) (err error)
	// Get TODO
	Get(userID uint64) (lifestyle *models.Lifestyle, err error)
}
