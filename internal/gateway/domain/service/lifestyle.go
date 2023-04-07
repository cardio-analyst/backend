package service

import (
	"database/sql"
	"errors"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether LifestyleService structure implements the service.LifestyleService interface
var _ service.LifestyleService = (*LifestyleService)(nil)

// LifestyleService implements service.LifestyleService interface.
type LifestyleService struct {
	lifestyles storage.LifestyleRepository
}

func NewLifestyleService(lifestyle storage.LifestyleRepository) *LifestyleService {
	return &LifestyleService{
		lifestyles: lifestyle,
	}
}

func (s *LifestyleService) Get(userID uint64) (*domain.Lifestyle, error) {
	lifestyle, err := s.lifestyles.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLifestyleNotFound
		}
		return nil, err
	}

	return lifestyle, nil
}

func (s *LifestyleService) Update(lifestyleData domain.Lifestyle) error {
	return s.lifestyles.Update(lifestyleData)
}
