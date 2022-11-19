package service

import (
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.ScoreService = (*scoreService)(nil)

// diseasesService implements service.DiseasesService interface.
type scoreService struct {
	cveRisk storage.ScoreRepository
}

func NewScoreService(cveRisk storage.ScoreRepository) *scoreService {
	return &scoreService{
		cveRisk: cveRisk,
	}
}

func (s scoreService) GetCveRisk(cveRiskData models.CveRiskData) (cveRisk uint64, err error) {
	riskValue, err := s.cveRisk.GetCveRisk(cveRiskData)

	if err != nil {
		return 0, err
	}

	return riskValue, nil
}
