package service

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const (
	diseasesNameTypeTwoDiabetes      = "Сахарный диабет 2 типа"
	diseasesNameInfarctionOrStroke   = "Перенесен инфаркт или инсульт"
	diseasesNameIschemicHeart        = "Ишемическая болезнь сердца"
	diseasesNameChronicKidney        = "Хроническая болезнь почек"
	diseasesNameArterialHypertension = "Артериальная гипертония"
	diseasesNameAtherosclerosis      = "Атеросклероз"
	diseasesNameOther                = "Другие СС-заболевания"
)

const (
	sbpRange80to114  = "80-114 мм.рт.ст"
	sbpRange115to149 = "115-149 мм.рт.ст"
	sbpRange150to184 = "150-184 мм.рт.ст"
	sbpRange185to219 = "185-219 мм.рт.ст"
	sbpRange220to250 = "220-250 мм.рт.ст"
)

type StatisticsService struct {
	registrationPublisher client.Publisher

	analyticsClient client.Analytics
	authClient      client.Auth

	diseases        storage.DiseasesRepository
	basicIndicators storage.BasicIndicatorsRepository
}

func NewStatisticsService(
	registrationPublisher client.Publisher,
	analyticsClient client.Analytics,
	authClient client.Auth,
	diseases storage.DiseasesRepository,
	basicIndicators storage.BasicIndicatorsRepository,
) *StatisticsService {
	return &StatisticsService{
		registrationPublisher: registrationPublisher,
		analyticsClient:       analyticsClient,
		authClient:            authClient,
		diseases:              diseases,
		basicIndicators:       basicIndicators,
	}
}

func (s *StatisticsService) NotifyRegistration(user model.User) {
	message := &model.MessageRegistration{
		Region: user.Region,
	}

	rmqMessageRaw, err := json.Marshal(message)
	if err != nil {
		log.Errorf("serializing registration message: %v", err)
		return
	}

	if err = s.registrationPublisher.Publish(rmqMessageRaw); err != nil {
		log.Errorf("publishing registration message: %v", err)
	}
}

func (s *StatisticsService) UsersByRegions(ctx context.Context) (map[string]int64, error) {
	return s.analyticsClient.UsersByRegions(ctx)
}

func (s *StatisticsService) DiseasesByUsers(region string, regionUsers map[uint64]bool) (map[string]int64, error) {
	diseases, err := s.diseases.All()
	if err != nil {
		return nil, err
	}

	result := map[string]int64{
		diseasesNameTypeTwoDiabetes:      0,
		diseasesNameInfarctionOrStroke:   0,
		diseasesNameIschemicHeart:        0,
		diseasesNameChronicKidney:        0,
		diseasesNameArterialHypertension: 0,
		diseasesNameAtherosclerosis:      0,
		diseasesNameOther:                0,
	}

	for _, disease := range diseases {
		if region != "" && !regionUsers[disease.UserID] {
			continue
		}

		if disease.HasChronicKidneyDisease {
			result[diseasesNameChronicKidney] += 1
		}
		if disease.HasArterialHypertension {
			result[diseasesNameArterialHypertension] += 1
		}
		if disease.HasIschemicHeartDisease {
			result[diseasesNameIschemicHeart] += 1
		}
		if disease.HasTypeTwoDiabetes {
			result[diseasesNameTypeTwoDiabetes] += 1
		}
		if disease.HadInfarctionOrStroke {
			result[diseasesNameInfarctionOrStroke] += 1
		}
		if disease.HasAtherosclerosis {
			result[diseasesNameAtherosclerosis] += 1
		}
		if disease.HasOtherCVD {
			result[diseasesNameOther] += 1
		}
	}

	return result, nil
}

func (s *StatisticsService) SBPByUsers(region string, regionUsers map[uint64]bool) (map[string]int64, error) {
	basicIndicators, err := s.basicIndicators.All()
	if err != nil {
		return nil, err
	}

	result := map[string]int64{
		sbpRange80to114:  0,
		sbpRange115to149: 0,
		sbpRange150to184: 0,
		sbpRange185to219: 0,
		sbpRange220to250: 0,
	}

	for _, basicIndicator := range basicIndicators {
		if basicIndicator.SBPLevel == nil {
			continue
		}

		if region != "" && !regionUsers[basicIndicator.UserID] {
			continue
		}

		switch {
		case 80.0 <= *basicIndicator.SBPLevel && *basicIndicator.SBPLevel < 115.0:
			result[sbpRange80to114] += 1
		case 115.0 <= *basicIndicator.SBPLevel && *basicIndicator.SBPLevel < 150.0:
			result[sbpRange115to149] += 1
		case 150.0 <= *basicIndicator.SBPLevel && *basicIndicator.SBPLevel < 185.0:
			result[sbpRange150to184] += 1
		case 185.0 <= *basicIndicator.SBPLevel && *basicIndicator.SBPLevel < 220.0:
			result[sbpRange185to219] += 1
		case 220.0 <= *basicIndicator.SBPLevel && *basicIndicator.SBPLevel < 251.0:
			result[sbpRange220to250] += 1
		}
	}

	return result, nil
}

func (s *StatisticsService) IdealCardiovascularAgesRangesByUsers(region string, regionUsers map[uint64]bool) (map[string]int64, error) {
	basicIndicators, err := s.basicIndicators.All()
	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)

	for _, basicIndicator := range basicIndicators {
		if basicIndicator.IdealCardiovascularAgesRange == nil {
			continue
		}

		if region != "" && !regionUsers[basicIndicator.UserID] {
			continue
		}

		_, found := result[*basicIndicator.IdealCardiovascularAgesRange]
		if !found {
			result[*basicIndicator.IdealCardiovascularAgesRange] = 0
		}
		result[*basicIndicator.IdealCardiovascularAgesRange] += 1
	}

	return result, nil
}
