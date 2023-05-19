package service

import (
	"github.com/cardio-analyst/backend/internal/analytics/ports/client"
	"github.com/cardio-analyst/backend/internal/analytics/ports/service"
	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
)

type Services struct {
	storage storage.Storage

	feedbackConsumer     client.Consumer
	registrationConsumer client.Consumer

	feedbackService   service.FeedbackService
	statisticsService service.StatisticsService
}

func NewServices(storage storage.Storage, feedbackConsumer, registrationConsumer client.Consumer) *Services {
	return &Services{
		storage:              storage,
		feedbackConsumer:     feedbackConsumer,
		registrationConsumer: registrationConsumer,
	}
}

func (s *Services) Feedback() service.FeedbackService {
	if s.feedbackService != nil {
		return s.feedbackService
	}

	s.feedbackService = NewFeedbackService(s.storage.Feedback(), s.feedbackConsumer)

	return s.feedbackService
}

func (s *Services) Statistics() service.StatisticsService {
	if s.statisticsService != nil {
		return s.statisticsService
	}

	s.statisticsService = NewStatisticsService(s.storage.RegionUsers(), s.registrationConsumer)

	return s.statisticsService
}
