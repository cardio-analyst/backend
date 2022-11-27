package service

import (
	"bytes"
	"fmt"
	"math"
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
	templateHealthyEating    = "healthy_eating.tmpl"
	templateSmoking          = "smoking.tmpl"
	templateSBPLevel         = "sbp_level.tmpl"
	templateCholesterolLevel = "cholesterol_level.tmpl"
	templateBMI              = "bmi.tmpl"
	templateLifestyle        = "lifestyle.tmpl"
)

// recommendation titles
const (
	titleHealthyEating    = "Здоровое питание"
	titleSmoking          = "Отказ от курения"
	titleSBPLevel         = "АД"
	titleCholesterolLevel = "Холестерин"
	titleBMI              = "Ожирение"
	titleLifestyle        = "Образ жизни"
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

	basicIndicators, err := s.basicIndicators.FindAll(userID)
	if err != nil {
		return nil, err
	}

	recommendation, err = s.getSmokingRecommendation(userID, basicIndicators)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.getSBPLevelRecommendation(basicIndicators[0])
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.getBMIRecommendation(basicIndicators)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.getLifestyleRecommendation()
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

func (s *recommendationsService) getSmokingRecommendation(userID uint64, basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {

	scoreData := getScoreData(basicIndicators)

	if !scoreData.Smoking {
		return nil, nil
	}

	var user *models.User
	user, err := s.users.GetByCriteria(models.UserCriteria{
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

func (s *recommendationsService) getSBPLevelRecommendation(basicIndicator *models.BasicIndicators) (*models.Recommendation, error) {
	sbpLevel := *basicIndicator.SBPLevel

	if sbpLevel >= 140 {
		tmplBuffer := &bytes.Buffer{}
		if err := s.templates.ExecuteTemplate(tmplBuffer, templateSBPLevel, map[string]interface{}{}); err != nil {
			return nil, err
		}

		return &models.Recommendation{
			Title:       titleSBPLevel,
			Description: tmplBuffer.String(),
		}, nil
	}

	return nil, nil
}

func (s *recommendationsService) getBMIRecommendation(basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {
	bmi := math.Pow(*basicIndicators[0].Weight/(*basicIndicators[0].Height/100), 2)
	scoreData := getScoreData(basicIndicators)

	var waist string

	if scoreData.Gender == "male" && bmi > 102 {
		waist = "Также у Вас превышен объем талии, необходимо уменьшить его минимум до 102"
	}

	if scoreData.Gender == "female" && bmi > 88 {
		waist = "Также у Вас превышен объем талии, необходимо уменьшить его минимум до 88"
	}

	tmplBuffer := &bytes.Buffer{}
	if err := s.templates.ExecuteTemplate(tmplBuffer, templateBMI, map[string]interface{}{
		"bmi":   fmt.Sprintf("%.2f", bmi),
		"waist": waist,
	}); err != nil {
		return nil, err
	}

	return &models.Recommendation{
		Title:       titleBMI,
		Description: tmplBuffer.String(),
	}, nil
}

func (s *recommendationsService) getLifestyleRecommendation() (*models.Recommendation, error) {
	tmplBuffer := &bytes.Buffer{}
	if err := s.templates.ExecuteTemplate(tmplBuffer, templateLifestyle, map[string]interface{}{}); err != nil {
		return nil, err
	}

	return &models.Recommendation{
		Title:       titleLifestyle,
		Description: tmplBuffer.String(),
	}, nil
}

func (s *recommendationsService) getCholesterolLevelRecommendation(basicIndicators []*models.BasicIndicators, userID uint64) (*models.Recommendation, error) {
	scoreData := getScoreData(basicIndicators)

	totalCholesterolLevel := scoreData.TotalCholesterolLevel
	gender := scoreData.Gender

	userDiseases, err := s.diseases.Get(userID)

	if err != nil {
		return nil, err
	}

	statusCholesterol := totalCholesterolLevel > 5 || (totalCholesterolLevel > 4.5 && userDiseases.HasTypeTwoDiabetes)

	if !statusCholesterol {
		return nil, nil
	}

	var minCholesterol string

	if userDiseases.HasTypeTwoDiabetes {
		minCholesterol = "4-4.5"
	} else {
		minCholesterol = "5"
	}

	var maxAlcohol string

	if gender == "female" {
		maxAlcohol = "10-20"
	} else {
		maxAlcohol = "20-30"
	}

	tmplBuffer := &bytes.Buffer{}
	if err = s.templates.ExecuteTemplate(tmplBuffer, templateCholesterolLevel, map[string]interface{}{
		"minCholesterol": minCholesterol,
		"maxAlcohol":     maxAlcohol,
	}); err != nil {
		return nil, err
	}

	return &models.Recommendation{
		Title:       titleCholesterolLevel,
		Description: tmplBuffer.String(),
	}, nil
}

func getScoreData(basicIndicators []*models.BasicIndicators) models.ScoreData {
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

	return scoreData
}
