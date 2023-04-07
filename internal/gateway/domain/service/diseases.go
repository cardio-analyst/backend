package service

import (
	"database/sql"
	"errors"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether DiseasesService structure implements the service.DiseasesService interface
var _ service.DiseasesService = (*DiseasesService)(nil)

// DiseasesService implements service.DiseasesService interface.
type DiseasesService struct {
	diseases storage.DiseasesRepository
}

func NewDiseasesService(diseases storage.DiseasesRepository) *DiseasesService {
	return &DiseasesService{
		diseases: diseases,
	}
}

func (s *DiseasesService) Get(userID uint64) (*domain.Diseases, error) {
	diseases, err := s.diseases.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserDiseasesNotFound
		}
		return nil, err
	}

	return diseases, nil
}

func (s *DiseasesService) Update(diseaseData domain.Diseases) error {
	return s.diseases.Update(diseaseData)
}
