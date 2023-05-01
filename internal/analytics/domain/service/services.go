package service

import (
	"github.com/cardio-analyst/backend/internal/analytics/ports/client"
	"github.com/cardio-analyst/backend/internal/analytics/ports/service"
	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
)

type Services struct {
	storage          storage.Storage
	feedbackConsumer client.FeedbackConsumer

	feedbackService service.FeedbackService
}

func NewServices(storage storage.Storage, feedbackConsumer client.FeedbackConsumer) *Services {
	return &Services{
		storage:          storage,
		feedbackConsumer: feedbackConsumer,
	}
}

func (s *Services) Feedback() service.FeedbackService {
	if s.feedbackService != nil {
		return s.feedbackService
	}

	s.feedbackService = NewFeedbackService(s.storage.Feedback(), s.feedbackConsumer)

	return s.feedbackService
}
