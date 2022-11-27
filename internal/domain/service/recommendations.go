package service

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	"github.com/cardio-analyst/backend/internal/domain/common"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// templates that contains recommendation contents pattern
const templatesPattern = "templates/recommendations/*.tmpl"

// recommendation template names
const (
	templateHealthyEating = "healthy_eating.tmpl"
	templateSmoking       = "smoking.tmpl"
)

// recommendation titles
const (
	titleHealthyEating = "Здоровое питание"
	titleSmoking       = "Отказ от курения"
)

// check whether recommendationsService structure implements the service.RecommendationsService interface
var _ service.RecommendationsService = (*recommendationsService)(nil)

// recommendationsService implements service.RecommendationsService interface.
type recommendationsService struct {
	diseases        storage.DiseasesRepository
	basicIndicators storage.BasicIndicatorsRepository
	lifestyles      storage.LifestyleRepository
	score           storage.ScoreRepository
	users           storage.UserRepository

	templates *template.Template
}

func NewRecommendationsService(
	diseases storage.DiseasesRepository,
	basicIndicators storage.BasicIndicatorsRepository,
	lifestyle storage.LifestyleRepository,
	score storage.ScoreRepository,
	users storage.UserRepository,
) *recommendationsService {
	return &recommendationsService{
		diseases:        diseases,
		basicIndicators: basicIndicators,
		lifestyles:      lifestyle,
		score:           score,
		users:           users,
		templates: template.Must(
			template.ParseGlob(templatesPattern),
		),
	}
}

func (s *recommendationsService) GetRecommendations(userID uint64) ([]*models.Recommendation, error) {
	recommendations := make([]*models.Recommendation, 0, 5)

	recommendation := s.getHealthyEatingRecommendation()
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	var err error
	recommendation, err = s.getSmokingRecommendation(userID)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

func (s *recommendationsService) getHealthyEatingRecommendation() *models.Recommendation {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100)%2 != 1 {
		return nil
	}

	tmplBuffer := &bytes.Buffer{}
	if err := s.templates.ExecuteTemplate(tmplBuffer, templateHealthyEating, map[string]interface{}{}); err != nil {
		return nil
	}

	return &models.Recommendation{
		Title:       titleHealthyEating,
		Description: tmplBuffer.String(),
	}
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

	tmplBuffer := &bytes.Buffer{}
	if err = s.templates.ExecuteTemplate(tmplBuffer, templateSmoking, map[string]interface{}{
		"riskSmoking":    riskSmoking,
		"riskNotSmoking": riskNotSmoking,
	}); err != nil {
		return nil, err
	}

	return &models.Recommendation{
		Title:       titleSmoking,
		Description: tmplBuffer.String(),
	}, nil
}
