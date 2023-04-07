package service

import (
	"database/sql"
	"errors"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether BasicIndicatorsService structure implements the service.BasicIndicatorsService interface
var _ service.BasicIndicatorsService = (*BasicIndicatorsService)(nil)

// BasicIndicatorsService implements service.BasicIndicatorsService interface.
type BasicIndicatorsService struct {
	basicIndicators storage.BasicIndicatorsRepository
}

func NewBasicIndicatorsService(basicIndicators storage.BasicIndicatorsRepository) *BasicIndicatorsService {
	return &BasicIndicatorsService{
		basicIndicators: basicIndicators,
	}
}

func (s *BasicIndicatorsService) Create(basicIndicatorsData domain.BasicIndicators) error {
	if err := basicIndicatorsData.Validate(false); err != nil {
		return err
	}

	return s.basicIndicators.Save(basicIndicatorsData)
}

func (s *BasicIndicatorsService) Update(basicIndicatorsData domain.BasicIndicators) error {
	if err := basicIndicatorsData.Validate(true); err != nil {
		return err
	}

	_, err := s.basicIndicators.Get(basicIndicatorsData.ID, basicIndicatorsData.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrBasicIndicatorsRecordNotFound
		}
		return err
	}

	return s.basicIndicators.Save(basicIndicatorsData)
}

func (s *BasicIndicatorsService) FindAll(userID uint64) ([]*domain.BasicIndicators, error) {
	return s.basicIndicators.FindAll(userID)
}
