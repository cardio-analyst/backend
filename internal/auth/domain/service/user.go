package service

import (
	"context"

	"github.com/alexedwards/argon2id"

	"github.com/cardio-analyst/backend/internal/auth/port/storage"
	"github.com/cardio-analyst/backend/pkg/model"
)

type UserService struct {
	users storage.UserRepository
}

func NewUserService(users storage.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Save(ctx context.Context, user model.User) error {
	criteria := model.UserCriteria{
		Login:             user.Login,
		Email:             user.Email,
		CriteriaSeparator: model.CriteriaSeparatorOR,
	}

	users, err := s.users.FindAllByCriteria(ctx, criteria)
	if err != nil {
		return err
	}

	// check if there are no users with a username or email that the actor wants to occupy
	if len(users) > 0 {
		for _, u := range users {
			if u.Login == user.Login {
				return model.ErrUserLoginAlreadyOccupied
			} else {
				return model.ErrUserEmailAlreadyOccupied
			}
		}
	}

	passwordHash, err := generateHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	return s.users.Save(ctx, user)
}

func generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *UserService) GetOne(ctx context.Context, criteria model.UserCriteria) (model.User, error) {
	return s.users.GetOneByCriteria(ctx, criteria)
}
