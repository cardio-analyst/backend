package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/pkg/model"
)

// check whether AuthService structure implements the service.AuthService interface
var _ service.AuthService = (*AuthService)(nil)

// AuthService implements service.AuthService interface.
type AuthService struct {
	client client.Auth
}

func NewAuthService(client client.Auth) *AuthService {
	return &AuthService{
		client: client,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user model.User) error {
	return s.client.SaveUser(ctx, user)
}

func (s *AuthService) GetTokens(ctx context.Context, credentials model.Credentials, userIP string, userRole model.UserRole) (model.Tokens, error) {
	return s.client.GetTokens(ctx, credentials, userIP, userRole)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken, userIP string, userRole model.UserRole) (model.Tokens, error) {
	return s.client.RefreshTokens(ctx, refreshToken, userIP, userRole)
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) (uint64, model.UserRole, error) {
	return s.client.IdentifyUser(ctx, token)
}

func (s *AuthService) GenerateSecretKey(ctx context.Context, userLogin, userEmail string) (string, error) {
	return s.client.GenerateSecretKey(ctx, userLogin, userEmail)
}
