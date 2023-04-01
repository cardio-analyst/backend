package service

import (
	models2 "github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// AuthService TODO
type AuthService interface {
	// RegisterUser TODO
	RegisterUser(user models2.User) (err error)
	// GetTokens TODO
	GetTokens(credentials models2.UserCredentials, userIP string) (tokens *models2.Tokens, err error)
	// ValidateAccessToken TODO
	ValidateAccessToken(token string) (userID uint64, err error)
	// RefreshTokens TODO
	RefreshTokens(refreshToken, userIP string) (tokens *models2.Tokens, err error)
}
