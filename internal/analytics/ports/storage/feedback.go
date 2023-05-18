package storage

import "github.com/cardio-analyst/backend/internal/pkg/model"

type FeedbackRepository interface {
	Create(feedback model.Feedback) error
	FindAll(criteria model.FeedbackCriteria) ([]model.Feedback, error)
	Count(criteria model.FeedbackCriteria) (int64, error)
	One(id uint64) (*model.Feedback, error)
	UpdateViewed(id uint64, viewed bool) error
}
