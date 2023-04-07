package service

import "github.com/cardio-analyst/backend/pkg/model"

type ValidationService interface {
	ValidateUser(user model.User, checkPassword bool) error
	ValidateDate(date model.Date) error
	ValidateCredentials(credentials model.Credentials) error
}
