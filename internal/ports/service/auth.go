package service

import "github.com/cardio-analyst/backend/internal/domain/models"

// AuthService TODO
type AuthService interface {
	// RegisterUser TODO
	RegisterUser(user models.User) (err error)
	// GetTokens TODO
	GetTokens(credentials models.UserCredentials, userIP string) (tokens *models.Tokens, err error)
	// ValidateAccessToken TODO
	ValidateAccessToken(token string) (userID uint64, err error)
	// RefreshTokens TODO
	RefreshTokens(refreshToken, userIP string) (tokens *models.Tokens, err error)
}
