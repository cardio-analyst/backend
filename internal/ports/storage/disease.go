package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

type DiseaseStorage interface {
	SaveDisease(diseaseData models.Disease) (err error)

	GetDiseaseByUserId(userId uint64) (disease *models.Disease, err error)
}
