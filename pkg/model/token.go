package model

import "errors"

var (
	ErrTokenIsExpired = errors.New("token is expired")
	ErrWrongToken     = errors.New("wrong token")
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
