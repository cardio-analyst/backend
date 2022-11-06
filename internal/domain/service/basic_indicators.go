package service

import (
	"database/sql"
	"errors"
	"fmt"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// check whether basicIndicatorsService structure implements the service.BasicIndicatorsService interface
var _ service.BasicIndicatorsService = (*basicIndicatorsService)(nil)

// basicIndicatorsService implements service.BasicIndicatorsService interface.
type basicIndicatorsService struct {
	basicIndicators storage.BasicIndicatorsRepository
}

func NewBasicIndicatorsService(basicIndicators storage.BasicIndicatorsRepository) *basicIndicatorsService {
	return &basicIndicatorsService{
		basicIndicators: basicIndicators,
	}
}

func (s *basicIndicatorsService) Create(basicIndicatorsData models.BasicIndicators) error {
	if err := basicIndicatorsData.Validate(false); err != nil {
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidBasicIndicatorsData, err)
	}

	return s.basicIndicators.Save(basicIndicatorsData)
}

func (s *basicIndicatorsService) Update(basicIndicatorsData models.BasicIndicators) error {
	if err := basicIndicatorsData.Validate(true); err != nil {
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidBasicIndicatorsData, err)
	}

	_, err := s.basicIndicators.Get(basicIndicatorsData.ID, basicIndicatorsData.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return serviceErrors.ErrBasicIndicatorsRecordNotFound
		}
		return err
	}

	return s.basicIndicators.Save(basicIndicatorsData)
}

func (s *basicIndicatorsService) FindAll(userID uint64) ([]*models.BasicIndicators, error) {
	return s.basicIndicators.FindAll(userID)
}
