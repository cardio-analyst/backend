package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.UserService = (*userService)(nil)

type userService struct {
	users storage.UserStorage
}

func NewUserService(users storage.UserStorage) *userService {
	return &userService{
		users: users,
	}
}

func (s *userService) Get(criteria models.UserCriteria) (*models.User, error) {
	user, err := s.users.GetUserByCriteria(criteria)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrUserNotFound
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
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidUserData, err)
	}

	criteria := models.UserCriteria{
		Login:             &user.Login,
		Email:             &user.Email,
		CriteriaSeparator: models.CriteriaSeparatorOR,
	}

	users, err := s.users.FindUserByCriteria(criteria)
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
				return serviceErrors.ErrUserLoginAlreadyOccupied
			} else if u.Email == user.Email {
				return serviceErrors.ErrUserEmailAlreadyOccupied
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

	return s.users.SaveUser(user)
}

func (s *userService) generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
