package service

import (
	"database/sql"
	"errors"
	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.DiseaseService = (*diseaseService)(nil)

type diseaseService struct {
	diseases storage.DiseaseRepository
}

func NewDiseaseService(diseases storage.DiseaseRepository) *diseaseService {
	return &diseaseService{
		diseases: diseases,
	}
}

func (d diseaseService) Get(userId uint64) (*models.Disease, error) {
	disease, err := d.diseases.GetByUserId(userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrUserNotFound
		}
		return nil, err
	}

	return disease, nil
}

func (d diseaseService) Update(diseaseData models.Disease) error {

	return d.diseases.Save(diseaseData)
}
