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

// check whether analysisService structure implements the service.AnalysisService interface
var _ service.AnalysisService = (*analysisService)(nil)

// analysisService implements service.AnalysisService interface.
type analysisService struct {
	analyses storage.AnalysisRepository
}

func NewAnalysisService(analyses storage.AnalysisRepository) *analysisService {
	return &analysisService{
		analyses: analyses,
	}
}

func (s *analysisService) Create(analysisData models.Analysis) error {
	if err := analysisData.Validate(false); err != nil {
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidAnalysisData, err)
	}

	return s.analyses.Save(analysisData)
}

func (s *analysisService) Update(analysisData models.Analysis) error {
	if err := analysisData.Validate(true); err != nil {
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidAnalysisData, err)
	}

	_, err := s.analyses.Get(analysisData.ID, analysisData.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return serviceErrors.ErrAnalysisRecordNotFound
		}
		return err
	}

	return s.analyses.Save(analysisData)
}

func (s *analysisService) FindAll(userID uint64) ([]*models.Analysis, error) {
	return s.analyses.FindAll(userID)
}
