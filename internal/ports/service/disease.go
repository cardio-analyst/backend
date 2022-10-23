package service

import "github.com/cardio-analyst/backend/internal/domain/models"

type DiseaseService interface {
	Get(userId uint64) (disease *models.Disease, err error)

	Update(diseaseData models.Disease) (err error)
}
