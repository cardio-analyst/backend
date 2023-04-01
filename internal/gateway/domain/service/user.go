package service

import (
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	errors2 "github.com/cardio-analyst/backend/internal/gateway/domain/errors"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether userService structure implements the service.UserService interface
var _ service.UserService = (*userService)(nil)

// userService implements service.UserService interface.
type userService struct {
	users storage.UserRepository
}

func NewUserService(users storage.UserRepository) *userService {
	return &userService{
		users: users,
	}
}

func (s *userService) Get(criteria models.UserCriteria) (*models.User, error) {
	user, err := s.users.GetByCriteria(criteria)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrUserNotFound
		}
		return nil, err
	}

	// we don't show passwords to anyone
	user.Password = ""

	return user, nil
}

func (s *userService) Update(user models.User) error {
	var validatePassword bool
	if user.Password != "" {
		validatePassword = true
	}

	if err := user.Validate(validatePassword); err != nil {
		return err
	}

	criteria := models.UserCriteria{
		Login:             &user.Login,
		Email:             &user.Email,
		CriteriaSeparator: models.CriteriaSeparatorOR,
	}

	users, err := s.users.FindByCriteria(criteria)
	if err != nil {
		return err
	}

	// check if there are no users with a username or email that the current user wants to occupy
	if len(users) > 0 {
		for _, u := range users {
			if u.ID == user.ID {
				continue
			}

			if u.Login == user.Login {
				return errors2.ErrUserLoginAlreadyOccupied
			} else if u.Email == user.Email {
				return errors2.ErrUserEmailAlreadyOccupied
			}
		}
	}

	if user.Password != "" {
		var passwordHash string
		passwordHash, err = s.generateHash(user.Password)
		if err != nil {
			return err
		}

		user.Password = passwordHash
	} else {
		user.Password = users[0].Password
	}

	return s.users.Save(user)
}

func (s *userService) generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
