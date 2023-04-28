package storage

import "github.com/cardio-analyst/backend/internal/analytics/domain/model"

type FeedbackRepository interface {
	Create(feedback model.Feedback) error
}
