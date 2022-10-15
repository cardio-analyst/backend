package errors

import "errors"

var (
	ErrInvalidUserData          = errors.New("invalid user data")
	ErrUserLoginAlreadyOccupied = errors.New("user with such login is already registered")
	ErrUserEmailAlreadyOccupied = errors.New("user with such email is already registered")
)
