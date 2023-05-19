package service

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type StatisticsService struct {
	registrationPublisher client.Publisher
	analyticsClient       client.Analytics
}

func NewStatisticsService(analyticsClient client.Analytics, registrationPublisher client.Publisher) *StatisticsService {
	return &StatisticsService{
		analyticsClient:       analyticsClient,
		registrationPublisher: registrationPublisher,
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
