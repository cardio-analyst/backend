package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

type DiseaseRepository interface {
	Save(diseaseData models.Disease) (err error)

	GetByUserId(userId uint64) (disease *models.Disease, err error)
}
