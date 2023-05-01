package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user model.User) (err error)
	GetTokens(ctx context.Context, credentials model.Credentials, userIP string, userRole model.UserRole) (tokens model.Tokens, err error)
	ValidateAccessToken(ctx context.Context, token string) (userID uint64, userRole model.UserRole, err error)
	RefreshTokens(ctx context.Context, refreshToken, userIP string, userRole model.UserRole) (tokens model.Tokens, err error)
	GenerateSecretKey(ctx context.Context, userLogin, userEmail string) (secretKey string, err error)
}
