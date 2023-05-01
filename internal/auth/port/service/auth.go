package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type AuthService interface {
	GetTokens(ctx context.Context, credentials model.Credentials, userIP string, userRole model.UserRole) (tokens model.Tokens, err error)
	RefreshTokens(ctx context.Context, refreshToken, userIP string, userRole model.UserRole) (tokens model.Tokens, err error)

	IdentifyUser(ctx context.Context, token string) (userID uint64, userRole model.UserRole, err error)

	GenerateSecretKey(userLogin, userEmail string) (secretKey string, err error)
	VerifySecretKey(user model.User) (err error)
}
