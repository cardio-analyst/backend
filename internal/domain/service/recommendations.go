package service

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	"github.com/cardio-analyst/backend/internal/config"
	"github.com/cardio-analyst/backend/internal/domain/common"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// recommendation template names
const (
	templateNameSmoking = "smoking"
)

// check whether recommendationsService structure implements the service.RecommendationsService interface
var _ service.RecommendationsService = (*recommendationsService)(nil)

// recommendationsService implements service.RecommendationsService interface.
type recommendationsService struct {
	cfg config.RecommendationsConfig

	diseases        storage.DiseasesRepository
	basicIndicators storage.BasicIndicatorsRepository
	lifestyles      storage.LifestyleRepository
	score           storage.ScoreRepository
	users           storage.UserRepository
}

func NewRecommendationsService(
	config config.RecommendationsConfig,
	diseases storage.DiseasesRepository,
	basicIndicators storage.BasicIndicatorsRepository,
	lifestyle storage.LifestyleRepository,
	score storage.ScoreRepository,
	users storage.UserRepository,
) *recommendationsService {
	return &recommendationsService{
		cfg:             config,
		diseases:        diseases,
		basicIndicators: basicIndicators,
		lifestyles:      lifestyle,
		score:           score,
		users:           users,
	}
}

func (s *recommendationsService) GetRecommendations(userID uint64) ([]*models.Recommendation, error) {
	recommendations := make([]*models.Recommendation, 0, 6)

	recommendation, err := s.getSmokingRecommendation(userID)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.getLifestyleRecommendation(userID)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.getHealthyEatingRecommendation()
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

func (s *recommendationsService) getHealthyEatingRecommendation() (*models.Recommendation, error) {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100)%2 != 1 {
		return nil, nil
	}

	return &models.Recommendation{
		What: s.cfg.HealthyEating.What,
		Why:  s.cfg.HealthyEating.Why,
		How:  s.cfg.HealthyEating.How,
	}, nil
}

func (s *recommendationsService) getSmokingRecommendation(userID uint64) (*models.Recommendation, error) {
	basicIndicators, err := s.basicIndicators.FindAll(userID)
	if err != nil {
		return nil, err
	}

	var scoreData models.ScoreData
	for _, basicIndicator := range basicIndicators {
		if basicIndicator.Smoking != nil && *basicIndicator.Smoking {
			scoreData.Smoking = true
		}
		if basicIndicator.Gender != nil && scoreData.Gender == "" {
			scoreData.Gender = *basicIndicator.Gender
		}
		if basicIndicator.SBPLevel != nil && scoreData.SBPLevel == 0 {
			scoreData.SBPLevel = *basicIndicator.SBPLevel
		}
		if basicIndicator.TotalCholesterolLevel != nil && scoreData.TotalCholesterolLevel == 0 {
			scoreData.TotalCholesterolLevel = *basicIndicator.TotalCholesterolLevel
		}

		// fastest break condition
		if scoreData.Gender != "" && scoreData.SBPLevel != 0 && scoreData.TotalCholesterolLevel != 0 {
			break
		}
	}

	if !scoreData.Smoking {
		return nil, nil
	}

	var user *models.User
	user, err = s.users.GetByCriteria(models.UserCriteria{
		ID: &userID,
	})
	if err != nil {
		return nil, err
	}

	scoreData.Age = common.GetCurrentAge(user.BirthDate.Time)

	var riskSmoking uint64
	riskSmoking, err = s.score.GetCVERisk(scoreData)
	if err != nil {
		return nil, err
	}

	scoreData.Smoking = false

	var riskNotSmoking uint64
	riskNotSmoking, err = s.score.GetCVERisk(scoreData)
	if err != nil {
		return nil, err
	}

	why, err := textTemplateToString(templateNameSmoking, s.cfg.Smoking.Why, map[string]interface{}{
		"riskSmoking":    riskSmoking,
		"riskNotSmoking": riskNotSmoking,
	})
	if err != nil {
		return nil, err
	}

	return &models.Recommendation{
		What: s.cfg.Smoking.What,
		Why:  why,
		How:  s.cfg.Smoking.How,
	}, nil
}

func (s *recommendationsService) getLifestyleRecommendation(userID uint64) (*models.Recommendation, error) {
	lifestyle, err := s.lifestyles.Get(userID)
	if err != nil {
		return nil, err
	}

	if lifestyle.EventsParticipation == common.EventsParticipationNotFrequently ||
		lifestyle.PhysicalActivity == common.PhysicalActivityOneInWeek {
		return &models.Recommendation{
			What: s.cfg.Lifestyle.What,
			Why:  s.cfg.Lifestyle.Why,
			How:  s.cfg.Lifestyle.How,
		}, nil
	}

	return nil, nil
}

func textTemplateToString(tmplName, tmplText string, tmplData map[string]interface{}) (string, error) {
	tmpl, err := template.New(tmplName).Parse(tmplText)
	if err != nil {
		return "", err
	}

	tmplBuffer := &bytes.Buffer{}
	if err = tmpl.Execute(tmplBuffer, tmplData); err != nil {
		return "", err
	}

	return tmplBuffer.String(), nil
}
