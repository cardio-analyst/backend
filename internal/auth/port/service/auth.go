package service

import (
	"context"

	"github.com/cardio-analyst/backend/pkg/model"
)

type AuthService interface {
	GetTokens(ctx context.Context, credentials model.Credentials, userIP string) (tokens model.Tokens, err error)
	RefreshTokens(ctx context.Context, refreshToken, userIP string) (tokens model.Tokens, err error)
	IdentifyUser(ctx context.Context, token string) (userID uint64, userRole model.UserRole, err error)
}
