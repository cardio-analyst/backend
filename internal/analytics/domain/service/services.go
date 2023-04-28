package service

import (
	"github.com/cardio-analyst/backend/internal/analytics/ports/service"
	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
)

type Services struct {
	storage storage.Storage

	feedbackService service.FeedbackService
}

func NewServices(storage storage.Storage) *Services {
	return &Services{
		storage: storage,
	}
}

func (s *Services) Feedback() service.FeedbackService {
	if s.feedbackService != nil {
		return s.feedbackService
	}

	s.feedbackService = NewFeedbackService(s.storage.Feedback())

	return s.feedbackService
}
