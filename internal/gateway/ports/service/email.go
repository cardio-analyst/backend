package service

import "github.com/cardio-analyst/backend/internal/pkg/model"

type EmailService interface {
	SendReport(receivers []string, reportFilePath string, user model.User) error
}
