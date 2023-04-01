package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// DiseasesService TODO
type DiseasesService interface {
	// Update TODO
	Update(diseasesData models.Diseases) (err error)
	// Get TODO
	Get(userID uint64) (diseasesData *models.Diseases, err error)
}
