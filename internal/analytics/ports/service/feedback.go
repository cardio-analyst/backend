package service

import "github.com/cardio-analyst/backend/internal/pkg/model"

type FeedbackService interface {
	ListenToFeedbackMessages() error
	FindAll() ([]model.Feedback, error)
	ToggleFeedbackViewed(id uint64) error
}
