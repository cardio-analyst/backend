package disease

import (
	"database/sql"
	"errors"
	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.DiseaseStorage = (*diseaseService)(nil)

type diseaseService struct {
	diseases storage.DiseaseStorage
}

func NewDiseaseService(diseases storage.DiseaseStorage) *diseaseService {
	return &diseaseService{
		diseases: diseases,
	}
}

func (d diseaseService) Get(userId uint) (*models.Disease, error) {
	disease, err := d.diseases.GetDiseaseByUserId(userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrUserNotFound
		}
		return nil, err
	}

	return disease, nil
}

func (d diseaseService) Update(diseaseData models.Disease) error {

	return d.diseases.SaveDisease(diseaseData)
}
