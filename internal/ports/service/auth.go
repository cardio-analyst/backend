package service

import "github.com/cardio-analyst/backend/internal/domain/models"

// AuthService TODO
type AuthService interface {
	// RegisterUser TODO
	RegisterUser(user models.User) (err error)
	// GetToken TODO
	GetToken(credentials models.UserCredentials) (token string, err error)
	// ValidateToken TODO
	ValidateToken(token string) (userID uint64, err error)
}
