package service

import "github.com/cardio-analyst/backend/internal/pkg/model"

type FeedbackService interface {
	Send(mark int16, message, version string, user model.User) error
	FindAll(criteria model.FeedbackCriteria) (feedbacks []model.Feedback, totalPages int64, err error)
	ToggleFeedbackViewed(id uint64) error
}
