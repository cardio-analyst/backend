package model

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid user credentials")
	ErrWrongCredentials   = errors.New("wrong user credentials")
)

type Credentials struct {
	Login    string
	Email    string
	Password string
}
