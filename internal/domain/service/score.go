package service

import (
	"fmt"

	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// check whether scoreService structure implements the service.ScoreService interface
var _ service.ScoreService = (*scoreService)(nil)

// scoreService implements service.ScoreService interface.
type scoreService struct {
	score storage.ScoreRepository
}

func NewScoreService(score storage.ScoreRepository) *scoreService {
	return &scoreService{
		score: score,
	}
}

func (s *scoreService) GetCVERisk(data models.ScoreData) (uint64, error) {
	if err := data.Validate(); err != nil {
		return 0, err
	}

	return s.score.GetCVERisk(data)
}

func (s *scoreService) GetIdealAge(data models.ScoreData) (string, error) {
	// pass SCORE data validation because GetCVERisk meth has it
	riskValue, err := s.GetCVERisk(data)
	if err != nil {
		return "", err
	}

	min, max, err := s.score.GetIdealAge(riskValue, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v-%v", min, max), nil
}
