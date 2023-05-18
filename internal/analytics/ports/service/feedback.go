package service

import "github.com/cardio-analyst/backend/internal/pkg/model"

type FeedbackService interface {
	ListenToFeedbackMessages() error
	FindAll(criteria model.FeedbackCriteria) (feedbacks []model.Feedback, totalPages int64, err error)
	ToggleFeedbackViewed(id uint64) error
}
