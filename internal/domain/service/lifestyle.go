package service

import (
	"database/sql"
	"errors"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// check whether lifestyleService structure implements the service.LifestyleService interface
var _ service.LifestyleService = (*lifestyleService)(nil)

// lifestyleService implements service.LifestyleService interface.
type lifestyleService struct {
	lifestyles storage.LifestyleRepository
}

func NewLifestyleService(lifestyle storage.LifestyleRepository) *lifestyleService {
	return &lifestyleService{
		lifestyles: lifestyle,
	}
}

func (s *lifestyleService) Get(userID uint64) (*models.Lifestyle, error) {
	lifestyle, err := s.lifestyles.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrLifestyleNotFound
		}
		return nil, err
	}

	return lifestyle, nil
}

func (s *lifestyleService) Update(lifestyleData models.Lifestyle) error {
	return s.lifestyles.Update(lifestyleData)
}
