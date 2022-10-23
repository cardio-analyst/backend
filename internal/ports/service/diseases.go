package service

import "github.com/cardio-analyst/backend/internal/domain/models"

// DiseasesService TODO
type DiseasesService interface {
	// Update TODO
	Update(diseasesData models.Diseases) (err error)
	// Get TODO
	Get(userID uint64) (diseases *models.Diseases, err error)
}
