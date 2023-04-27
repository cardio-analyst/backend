package model

import "errors"

var (
	ErrTokenIsExpired = errors.New("token is expired")
	ErrWrongToken     = errors.New("wrong token")

	ErrInvalidSecretKey = errors.New("invalid secret key")
	ErrWrongSecretKey   = errors.New("wrong secret key")

	ErrForbiddenByRole = errors.New("forbidden by role")
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
