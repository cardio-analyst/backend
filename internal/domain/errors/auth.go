package errors

import "errors"

var (
	ErrInvalidUserCredentials = errors.New("invalid user credentials")
	ErrWrongCredentials       = errors.New("wrong user credentials")
	ErrTokenIsExpired         = errors.New("token is expired")
	ErrWrongToken             = errors.New("wrong token")
	ErrIPIsNotInWhitelist     = errors.New("user ip is not in the whitelist")
	ErrSessionNotFound        = errors.New("session not found")
)
