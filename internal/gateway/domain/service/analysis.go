package service

import (
	"database/sql"
	"errors"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether AnalysisService structure implements the service.AnalysisService interface
var _ service.AnalysisService = (*AnalysisService)(nil)

// AnalysisService implements service.AnalysisService interface.
type AnalysisService struct {
	analyses storage.AnalysisRepository
}

func NewAnalysisService(analyses storage.AnalysisRepository) *AnalysisService {
	return &AnalysisService{
		analyses: analyses,
	}
}

func (s *AnalysisService) Create(analysisData domain.Analysis) error {
	if err := analysisData.Validate(false); err != nil {
		return err
	}

	return s.analyses.Save(analysisData)
}

func (s *AnalysisService) Update(analysisData domain.Analysis) error {
	if err := analysisData.Validate(true); err != nil {
		return err
	}

	_, err := s.analyses.Get(analysisData.ID, analysisData.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrAnalysisRecordNotFound
		}
		return err
	}

	return s.analyses.Save(analysisData)
}

func (s *AnalysisService) FindAll(userID uint64) ([]*domain.Analysis, error) {
	return s.analyses.FindAll(userID)
}
