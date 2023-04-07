package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/pkg/model"
)

// check whether UserService structure implements the service.UserService interface
var _ service.UserService = (*UserService)(nil)

// UserService implements service.UserService interface.
type UserService struct {
	authClient client.Auth
}

func NewUserService(authClient client.Auth) *UserService {
	return &UserService{
		authClient: authClient,
	}
}

func (s *UserService) GetOne(ctx context.Context, criteria model.UserCriteria) (model.User, error) {
	user, err := s.authClient.GetUser(ctx, criteria)
	if err != nil {
		return model.User{}, err
	}

	// we don't show passwords to anyone
	user.Password = ""

	return user, nil
}

func (s *UserService) Update(ctx context.Context, user model.User) error {
	return s.authClient.SaveUser(ctx, user)
}
