package service

import (
	"fmt"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// check whether scoreService structure implements the service.ScoreService interface
var _ service.ScoreService = (*scoreService)(nil)

// scoreService implements service.ScoreService interface.
type scoreService struct {
	cveRisk storage.ScoreRepository
}

func NewScoreService(cveRisk storage.ScoreRepository) *scoreService {
	return &scoreService{
		cveRisk: cveRisk,
	}
}

func (s scoreService) GetCVERisk(cveRiskData models.CVERiskData) (uint64, error) {
	if err := cveRiskData.Validate(); err != nil {
		return 0, fmt.Errorf("%w: %v", serviceErrors.ErrInvalidCVERiskData, err)
	}

	return s.cveRisk.GetCVERisk(cveRiskData)
}
