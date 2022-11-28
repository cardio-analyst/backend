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
	templateNameBmi              = "bmi"
	templateNameCholesterolLevel = "cholesterol_level"
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

	var err error

	basicIndicators, err := s.basicIndicators.FindAll(userID)

	if err != nil {
		return nil, err
	}

	recommendation, err := s.getSmokingRecommendation(userID, basicIndicators)
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

	recommendation, err = s.getSBPLevelRecommendation(basicIndicators)
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

	recommendation, err = s.getCholesterolLevelRecommendation(userID, basicIndicators)
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

func (s *recommendationsService) getSBPLevelRecommendation(basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {
	var sbpLevel float64

	for _, basicIndicator := range basicIndicators {
		if basicIndicator.SBPLevel != nil && sbpLevel == 0 {
			sbpLevel = *basicIndicator.SBPLevel
			break
		}
	}

	if sbpLevel >= 140 {
		return &models.Recommendation{
			What: s.cfg.SbpLevel.What,
			Why:  s.cfg.SbpLevel.Why,
			How:  s.cfg.SbpLevel.How,
		}, nil
	}

	return nil, nil
}

func (s *recommendationsService) getBMIRecommendation(basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {

	var weight, height, waistSize float64

	for _, basicIndicator := range basicIndicators {
		if basicIndicator.Weight != nil && weight == 0 {
			weight = *basicIndicator.Weight
		}
		if basicIndicator.Height != nil && height == 0 {
			height = *basicIndicator.Height
		}
		if basicIndicator.WaistSize != nil && waistSize == 0 {
			waistSize = *basicIndicator.WaistSize
		}
		if weight != 0 && height != 0 && waistSize != 0 {
			break
		}
	}

	if weight == 0 || height == 0 {
		return nil, nil
	}

	bmi := weight / math.Pow(height/100, 2)
	scoreData := getScoreData(basicIndicators)

	var waistPrint string

	if scoreData.Gender == common.UserGenderMale && waistSize > 102 {
		waistPrint = " Также у Вас превышен объем талии, необходимо уменьшить его минимум до 102."
	}

	if scoreData.Gender == common.UserGenderFemale && waistSize > 88 {
		waistPrint = " Также у Вас превышен объем талии, необходимо уменьшить его минимум до 88."
	}

	why, err := textTemplateToString(templateNameBmi, s.cfg.Bmi.Why, map[string]interface{}{
		"bmi":   fmt.Sprintf("%.2f", bmi),
		"waist": waistPrint,
	})

	if err != nil {
		return nil, err
	}

	return &models.Recommendation{
		What: s.cfg.Bmi.What,
		Why:  why,
		How:  s.cfg.Bmi.How,
	}, nil
}

func (s *recommendationsService) getCholesterolLevelRecommendation(userID uint64, basicIndicators []*models.BasicIndicators) (*models.Recommendation, error) {
	scoreData := getScoreData(basicIndicators)

	if scoreData.TotalCholesterolLevel == 0 || scoreData.Gender == common.UserGenderUnknown {
		return nil, nil
	}

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

	why, err := textTemplateToString(templateNameCholesterolLevel, s.cfg.CholesterolLevel.Why, map[string]interface{}{
		"minCholesterol": minCholesterol,
	})

	if err != nil {
		return nil, err
	}

	bmi := math.Pow(*basicIndicators[0].Weight/(*basicIndicators[0].Height/100), 2)

	how, err := textTemplateToString(templateNameCholesterolLevel, s.cfg.CholesterolLevel.How, map[string]interface{}{
		"minSize":    fmt.Sprintf("%.2f", bmi),
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
