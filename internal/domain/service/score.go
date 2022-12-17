package service

import (
	"fmt"
	"github.com/cardio-analyst/backend/internal/domain/common"

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

func (s *scoreService) GetCVERisk(data models.ScoreData) (uint64, string, error) {
	if err := data.Validate(models.ValidationOptionsScore{
		Age:                   true,
		Gender:                true,
		SBPLevel:              true,
		TotalCholesterolLevel: true,
	}); err != nil {
		return 0, common.ScaleUnknown, err
	}

	riskValue, err := s.score.GetCVERisk(data)
	if err != nil {
		return 0, common.ScaleUnknown, err
	}

	scale := s.ResolveScale(float64(riskValue), data.Age)

	return riskValue, scale, nil
}

func (s *scoreService) ResolveScale(riskValue float64, age int) string {
	switch {
	case age > 0 && age < 50:
		switch {
		case riskValue > 0 && riskValue < 2.5:
			return common.ScalePositive
		case riskValue >= 2.5 && riskValue < 7.5:
			return common.ScaleNeutral
		case riskValue >= 7.5:
			return common.ScaleNegative
		default:
			return common.ScaleUnknown
		}
	case age >= 50 && age <= 69:
		switch {
		case riskValue > 0 && riskValue < 5.0:
			return common.ScalePositive
		case riskValue >= 5.0 && riskValue < 10.0:
			return common.ScaleNeutral
		case riskValue >= 10.0:
			return common.ScaleNegative
		default:
			return common.ScaleUnknown
		}
	case age >= 70:
		switch {
		case riskValue > 0 && riskValue < 7.5:
			return common.ScalePositive
		case riskValue >= 7.5 && riskValue < 15.0:
			return common.ScaleNeutral
		case riskValue >= 15.0:
			return common.ScaleNegative
		default:
			return common.ScaleUnknown
		}
	default:
		return common.ScaleUnknown
	}
}

func (s *scoreService) GetIdealAge(data models.ScoreData) (string, string, error) {
	// pass SCORE data validation because GetCVERisk meth has it
	riskValue, scale, err := s.GetCVERisk(data)
	if err != nil {
		return "", common.ScaleUnknown, err
	}

	min, max, err := s.score.GetIdealAge(riskValue, data)
	if err != nil {
		return "", common.ScaleUnknown, err
	}

	return fmt.Sprintf("%v-%v", min, max), scale, nil
}
