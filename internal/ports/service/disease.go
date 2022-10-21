package service

import "github.com/cardio-analyst/backend/internal/domain/models"

type DiseaseStorage interface {
	Get(userId uint) (disease *models.Disease, err error)

	Update(diseaseData models.Disease) (err error)
}
