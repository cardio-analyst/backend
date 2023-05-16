package service

import (
	"context"

	"github.com/alexedwards/argon2id"

	"github.com/cardio-analyst/backend/internal/auth/port/storage"
	"github.com/cardio-analyst/backend/internal/pkg/model"
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
			if u.ID == user.ID {
				continue
			}

			if u.Login == user.Login {
				return model.ErrUserLoginAlreadyOccupied
			} else if u.Email == user.Email {
				return model.ErrUserEmailAlreadyOccupied
			}
		}
	}

	if user.Password != "" {
		var passwordHash string
		passwordHash, err = generateHash(user.Password)
		if err != nil {
			return err
		}

		user.Password = passwordHash
	} else {
		user.Password = users[0].Password
	}

	return s.users.Save(ctx, user)
}

func (s *UserService) GetOne(ctx context.Context, criteria model.UserCriteria) (model.User, error) {
	user, err := s.users.GetOneByCriteria(ctx, criteria)
	if err != nil {
		return model.User{}, err
	}

	// we don't show password to anyone
	user.Password = ""

	return user, nil
}

func (s *UserService) GetList(ctx context.Context, criteria model.UserCriteria) ([]model.User, bool, error) {
	users, err := s.users.FindAllByCriteria(ctx, criteria)
	if err != nil {
		return nil, false, err
	}

	var hasNextPage bool
	if criteria.Limit > 0 {
		if len(users) == int(criteria.Limit+1) {
			users = users[:criteria.Limit]
			hasNextPage = true
		}
	}

	return users, hasNextPage, nil
}

func generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
