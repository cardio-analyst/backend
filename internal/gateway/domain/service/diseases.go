package service

import (
	"database/sql"
	"errors"
	serviceErrors "github.com/cardio-analyst/backend/internal/gateway/domain/errors"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether diseasesService structure implements the service.DiseasesService interface
var _ service.DiseasesService = (*diseasesService)(nil)

// diseasesService implements service.DiseasesService interface.
type diseasesService struct {
	diseases storage.DiseasesRepository
}

func NewDiseasesService(diseases storage.DiseasesRepository) *diseasesService {
	return &diseasesService{
		diseases: diseases,
	}
}

func (s *diseasesService) Get(userID uint64) (*models.Diseases, error) {
	diseases, err := s.diseases.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrUserDiseasesNotFound
		}
		return nil, err
	}

	return diseases, nil
}

func (s *diseasesService) Update(diseaseData models.Diseases) error {
	return s.diseases.Update(diseaseData)
}
