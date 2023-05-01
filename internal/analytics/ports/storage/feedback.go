package storage

import "github.com/cardio-analyst/backend/internal/pkg/model"

type FeedbackRepository interface {
	Create(feedback model.Feedback) error
	FindAll() ([]model.Feedback, error)
}
