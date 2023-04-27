package client

import (
	"context"

	"github.com/cardio-analyst/backend/pkg/model"
)

type Auth interface {
	SaveUser(ctx context.Context, user model.User) (err error)
	GetUser(ctx context.Context, criteria model.UserCriteria) (user model.User, err error)
	IdentifyUser(ctx context.Context, token string) (userID uint64, userRole model.UserRole, err error)

	GetTokens(ctx context.Context, credentials model.Credentials, userIP string, userRole model.UserRole) (tokens model.Tokens, err error)
	RefreshTokens(ctx context.Context, refreshToken, userIP string, userRole model.UserRole) (tokens model.Tokens, err error)

	GenerateSecretKey(ctx context.Context, userLogin, userEmail string) (secretKey string, err error)
}
