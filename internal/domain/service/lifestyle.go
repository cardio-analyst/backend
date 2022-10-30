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
	lifestyle storage.LifestyleRepository
}

func NewLifestyleService(lifestyle storage.LifestyleRepository) *lifestyleService {
	return &lifestyleService{
		lifestyle: lifestyle,
	}
}

func (l lifestyleService) Get(userID uint64) (*models.Lifestyle, error) {
	lifestyle, err := l.lifestyle.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrUserDiseasesNotFound
		}
		return nil, err
	}

	return lifestyle, nil
}

func (l lifestyleService) Update(lifestyleData models.Lifestyle) (err error) {
	return l.lifestyle.Update(lifestyleData)
}
