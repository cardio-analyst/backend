package service

import "github.com/cardio-analyst/backend/internal/domain/models"

type EmailService interface {
	SendReport(receivers []string, reportPath string, userData models.User) error
}
