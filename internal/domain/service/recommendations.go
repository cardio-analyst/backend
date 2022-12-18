package service

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"text/template"
	"time"

	"github.com/cardio-analyst/backend/internal/config"
	"github.com/cardio-analyst/backend/internal/domain/common"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// recommendation titles
const (
	templateNameSmoking          = "smoking"
	templateNameBMI              = "bmi"
	templateNameCholesterolLevel = "cholesterol_level"
	templateNameRisk             = "risk"
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
	cfg config.RecommendationsConfig,
	diseases storage.DiseasesRepository,
	basicIndicators storage.BasicIndicatorsRepository,
	lifestyle storage.LifestyleRepository,
	score storage.ScoreRepository,
	users storage.UserRepository,
) *recommendationsService {
	return &recommendationsService{
		cfg:             cfg,
		diseases:        diseases,
		basicIndicators: basicIndicators,
		lifestyles:      lifestyle,
		score:           score,
		users:           users,
	}
}

func (s *recommendationsService) GetRecommendations(userID uint64) ([]*models.Recommendation, error) {
	recommendations := make([]*models.Recommendation, 0, 7)

	basicIndicators, err := s.basicIndicators.FindAll(userID)
	if err != nil {
		return nil, err
	}

	recommendation, err := s.lifestyleRecommendation(userID)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation = s.healthyEatingRecommendation()
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	scoreData := models.ExtractScoreDataFrom(basicIndicators)

	recommendation, err = s.smokingRecommendation(userID, scoreData)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.sbpLevelRecommendation(scoreData)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.bmiRecommendation(scoreData, basicIndicators)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.cholesterolLevelRecommendation(userID, scoreData, basicIndicators)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	recommendation, err = s.riskRecommendation(userID, scoreData)
	if err != nil {
		return nil, err
	}
	if recommendation != nil {
		recommendations = append(recommendations, recommendation)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(recommendations), func(i, j int) {
		recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
	})

	return recommendations, nil
}

func (s *recommendationsService) healthyEatingRecommendation() *models.Recommendation {
	return &models.Recommendation{
		What: s.cfg.HealthyEating.What,
		Why:  s.cfg.HealthyEating.Why,
		How:  s.cfg.HealthyEating.How,
	}
}

func (s *recommendationsService) smokingRecommendation(userID uint64, scoreData models.ScoreData) (*models.Recommendation, error) {
	if err := scoreData.ValidateByRecommendation(models.Smoking); err != nil {
		return nil, nil
	}

	if !scoreData.Smoking {
		return nil, nil
	}

	user, err := s.users.GetByCriteria(models.UserCriteria{
		ID: &userID,
	})
	if err != nil {
		return nil, err
	}

	scoreData.Age = user.Age()

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

func (s *recommendationsService) sbpLevelRecommendation(scoreData models.ScoreData) (*models.Recommendation, error) {
	if err := scoreData.ValidateByRecommendation(models.SBPLevel); err != nil {
		return nil, nil
	}

	if scoreData.SBPLevel >= 140 {
		return &models.Recommendation{
			What: s.cfg.SBPLevel.What,
			Why:  s.cfg.SBPLevel.Why,
			How:  s.cfg.SBPLevel.How,
		}, nil
	}

	return nil, nil
}

func (s *recommendationsService) bmiRecommendation(scoreData models.ScoreData, basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {
	if err := scoreData.ValidateByRecommendation(models.BMI); err != nil {
		return nil, nil
	}

	weight, height, waistSize, bodyMassIndex := extractBMIIndications(basicIndicators)
	if bodyMassIndex < 25 {
		if weight == 0 || height == 0 {
			return nil, nil
		}

		bodyMassIndex = weight / math.Pow(height/100, 2)
		if bodyMassIndex < 25 {
			return nil, nil
		}
	}

	var waistRecommendation string
	switch {
	case scoreData.Gender == common.UserGenderMale && waistSize > 102:
		waistRecommendation = " Также у Вас превышен объем талии, необходимо уменьшить его минимум до 102."
	case scoreData.Gender == common.UserGenderFemale && waistSize > 88:
		waistRecommendation = " Также у Вас превышен объем талии, необходимо уменьшить его минимум до 88."
	}

	why, err := textTemplateToString(templateNameBMI, s.cfg.BMI.Why, map[string]interface{}{
		"bmi":   fmt.Sprintf("%.2f", bodyMassIndex),
		"waist": waistRecommendation,
	})
	if err != nil {
		return nil, err
	}

	return &models.Recommendation{
		What: s.cfg.BMI.What,
		Why:  why,
		How:  s.cfg.BMI.How,
	}, nil
}

func (s *recommendationsService) cholesterolLevelRecommendation(userID uint64, scoreData models.ScoreData, basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {
	if err := scoreData.ValidateByRecommendation(models.CholesterolLevel); err != nil {
		return nil, nil
	}

	if scoreData.TotalCholesterolLevel == 0 || scoreData.Gender == common.UserGenderUnknown {
		return nil, nil
	}

	userDiseases, err := s.diseases.Get(userID)
	if err != nil {
		return nil, err
	}

	statusCholesterol := scoreData.TotalCholesterolLevel > 5 || (scoreData.TotalCholesterolLevel > 4.5 && userDiseases.HasTypeTwoDiabetes)
	if !statusCholesterol {
		return nil, nil
	}

	var minCholesterol string
	if userDiseases.HasTypeTwoDiabetes {
		minCholesterol = "4-4.5"
	} else {
		minCholesterol = "5"
	}

	why, err := textTemplateToString(templateNameCholesterolLevel, s.cfg.CholesterolLevel.Why, map[string]interface{}{
		"minCholesterol": minCholesterol,
	})
	if err != nil {
		return nil, err
	}

	var maxAlcohol string
	if scoreData.Gender == common.UserGenderFemale {
		maxAlcohol = "10-20"
	} else {
		maxAlcohol = "20-30"
	}

	var weight, height, bodyMassIndex float64
	for _, indicators := range basicIndicators {
		if indicators.Weight != nil && weight == 0 {
			weight = *indicators.Weight
		}
		if indicators.Height != nil && height == 0 {
			height = *indicators.Height
		}
		if indicators.BodyMassIndex != nil && bodyMassIndex == 0 {
			bodyMassIndex = *indicators.BodyMassIndex
		}

		// fastest break condition
		if weight != 0 && height != 0 && bodyMassIndex != 0 {
			break
		}
	}

	if bodyMassIndex == 0 {
		bodyMassIndex = weight / math.Pow(height/100, 2)
		if bodyMassIndex == 0 {
			return nil, nil
		}
	}

	how, err := textTemplateToString(templateNameCholesterolLevel, s.cfg.CholesterolLevel.How, map[string]interface{}{
		"minSize":    fmt.Sprintf("%.2f", bodyMassIndex),
		"maxAlcohol": maxAlcohol,
	})
	if err != nil {
		return nil, err
	}

	return &models.Recommendation{
		What: s.cfg.CholesterolLevel.What,
		Why:  why,
		How:  how,
	}, nil
}

func (s *recommendationsService) lifestyleRecommendation(userID uint64) (*models.Recommendation, error) {
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

func (s *recommendationsService) riskRecommendation(userID uint64, scoreData models.ScoreData) (*models.Recommendation, error) {
	if err := scoreData.ValidateByRecommendation(models.Risk); err != nil {
		return nil, nil
	}

	user, err := s.users.GetByCriteria(models.UserCriteria{
		ID: &userID,
	})
	if err != nil {
		return nil, err
	}

	scoreData.Age = user.Age()

	riskActual, err := s.score.GetCVERisk(scoreData)
	if err != nil {
		return nil, err
	}

	ageMinActual, ageMaxActual, err := s.score.GetIdealAge(riskActual, scoreData)
	if err != nil {
		return nil, err
	}

	ageMinDifference := int(ageMinActual) - scoreData.Age
	ageMaxDifference := int(ageMaxActual) - scoreData.Age
	if ageMinDifference <= 0 || ageMaxDifference <= 0 {
		return nil, nil
	}

	why, err := textTemplateToString(templateNameRisk, s.cfg.Risk.Why, map[string]interface{}{
		"riskActual":          riskActual,
		"agesRangeActual":     fmt.Sprintf("%v-%v", ageMinActual, ageMaxActual),
		"agesRangeDifference": fmt.Sprintf("%v-%v", ageMinDifference, ageMaxDifference),
	})
	if err != nil {
		return nil, err
	}

	return &models.Recommendation{
		What: s.cfg.Risk.What,
		Why:  why,
		How:  s.cfg.Risk.How,
	}, nil
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

func extractBMIIndications(basicIndicators []*models.BasicIndicators) (weight, height, waistSize, bodyMassIndex float64) {
	for _, indicators := range basicIndicators {
		if indicators.Weight != nil && weight == 0 {
			weight = *indicators.Weight
		}
		if indicators.Height != nil && height == 0 {
			height = *indicators.Height
		}
		if indicators.WaistSize != nil && waistSize == 0 {
			waistSize = *indicators.WaistSize
		}
		if indicators.BodyMassIndex != nil && bodyMassIndex == 0 {
			bodyMassIndex = *indicators.BodyMassIndex
		}

		// fastest break condition
		if weight != 0 && height != 0 && waistSize != 0 && bodyMassIndex != 0 {
			break
		}
	}
	return
}
