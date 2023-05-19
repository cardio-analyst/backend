package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/cardio-analyst/backend/internal/analytics/ports/client"
	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type StatisticsService struct {
	regionUsersRepository storage.RegionUsersRepository
	registrationConsumer  client.Consumer
}

func NewStatisticsService(repository storage.RegionUsersRepository, consumer client.Consumer) *StatisticsService {
	return &StatisticsService{
		regionUsersRepository: repository,
		registrationConsumer:  consumer,
	}
}

func (s *StatisticsService) ListenToRegistrationMessages() error {
	handler := s.registrationMessagesHandler()

	if err := s.registrationConsumer.Consume(handler); err != nil {
		return fmt.Errorf("consuming registration messages: %w", err)
	}

	return nil
}

func (s *StatisticsService) registrationMessagesHandler() func(data []byte) error {
	return func(data []byte) error {
		var message model.MessageRegistration
		if err := json.Unmarshal(data, &message); err != nil {
			log.Errorf("unmarshalling registration message body: %v", err)
			return err
		}

		if err := s.regionUsersRepository.Increment(message.Region); err != nil {
			log.Errorf("incrementing users number for region %q: %v", message.Region, err)
			return err
		}

		return nil
	}
}

func (s *StatisticsService) AllUsersByRegions() (map[string]int64, error) {
	return s.regionUsersRepository.All()
}
