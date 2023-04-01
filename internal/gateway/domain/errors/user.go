package errors

import "errors"

var (
	ErrInvalidFirstName = errors.New("invalid firstName value")
	ErrInvalidLastName  = errors.New("invalid lastName value")
	ErrInvalidRegion    = errors.New("invalid region value")
	ErrInvalidBirthDate = errors.New("invalid birthDate value")
	ErrInvalidLogin     = errors.New("invalid login value")
	ErrInvalidEmail     = errors.New("invalid email value")
	ErrInvalidPassword  = errors.New("invalid password value")
	ErrUserNotFound     = errors.New("user not found")
)
