package service

import "github.com/cardio-analyst/backend/pkg/model"

type ValidationService interface {
	ValidateUser(user model.User) error
	ValidateDate(date model.Date) error
	ValidateCredentials(credentials model.Credentials) error
}
