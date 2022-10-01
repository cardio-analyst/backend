package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"

	"github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	db storage.UserStorage
}

func NewAuthService(db storage.UserStorage) *authService {
	return &authService{db: db}
}

func (s *authService) RegisterUser(user models.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidUserData, err)
	}

	found, err := s.db.FindOneByCriteria(models.UserCriteria{
		Login:             &user.Login,
		Email:             &user.Email,
		CriteriaSeparator: models.CriteriaSeparatorOR,
	})
	if err != nil {
		return err
	}
	if found != nil {
		if found.Login == user.Login {
			return errors.ErrUserLoginAlreadyOccupied
		} else {
			return errors.ErrUserEmailAlreadyOccupied
		}
	}

	passwordHash, err := s.generatePasswordHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	return s.db.Create(user)
}

func (s *authService) GetToken(credentials models.Credentials) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *authService) generatePasswordHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
