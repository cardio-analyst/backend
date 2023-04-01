package storage

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// LifestyleRepository TODO
type LifestyleRepository interface {
	// Update TODO
	Update(lifestyleData models.Lifestyle) (err error)
	// Get TODO
	Get(userID uint64) (lifestyleData *models.Lifestyle, err error)
}
