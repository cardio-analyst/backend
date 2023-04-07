package service

import (
	"fmt"

	"github.com/cardio-analyst/backend/internal/gateway/domain/common"
	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether ScoreService structure implements the service.ScoreService interface
var _ service.ScoreService = (*ScoreService)(nil)

// ScoreService implements service.ScoreService interface.
type ScoreService struct {
	score storage.ScoreRepository
}

func NewScoreService(score storage.ScoreRepository) *ScoreService {
	return &ScoreService{
		score: score,
	}
}

func (s *ScoreService) GetCVERisk(data domain.ScoreData) (uint64, string, error) {
	if err := data.Validate(domain.ValidationOptionsScore{
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

func (s *ScoreService) ResolveScale(riskValue float64, age int) string {
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

func (s *ScoreService) GetIdealAge(data domain.ScoreData) (string, string, error) {
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
