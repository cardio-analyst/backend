package service

import "github.com/cardio-analyst/backend/internal/pkg/model"

type FeedbackService interface {
	Send(mark int16, message, version string, user model.User) error
	FindAll() ([]model.Feedback, error)
}
