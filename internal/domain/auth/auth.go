package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"

	"github.com/cardio-analyst/backend/internal/config"
	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	cfg config.AuthConfig
	db  storage.UserStorage
}

func NewAuthService(cfg config.AuthConfig, db storage.UserStorage) *authService {
	return &authService{
		cfg: cfg,
		db:  db,
	}
}

func (s *authService) RegisterUser(user models.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidUserData, err)
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
			return serviceErrors.ErrUserLoginAlreadyOccupied
		} else {
			return serviceErrors.ErrUserEmailAlreadyOccupied
		}
	}

	passwordHash, err := s.generateHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	return s.db.Create(user)
}

func (s *authService) GetToken(credentials models.UserCredentials) (string, error) {
	if err := credentials.Validate(); err != nil {
		return "", fmt.Errorf("%w: %v", serviceErrors.ErrInvalidUserCredentials, err)
	}

	// since we do not know what exactly we are dealing with, we are looking for two fields
	criteria := models.UserCriteria{
		Login: &credentials.LoginOrEmail,
		Email: &credentials.LoginOrEmail,
	}

	user, err := s.db.GetOneByCriteria(criteria)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", serviceErrors.ErrWrongCredentials
		}
		return "", err
	}

	match, err := s.comparePasswordAndHash(credentials.Password, user.Password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", serviceErrors.ErrWrongCredentials
	}

	return s.generateJWTToken(user.Login)
}

type tokenClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func (s *authService) generateJWTToken(login string) (string, error) {
	claims := tokenClaims{
		login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.cfg.SigningKey)
}

func (s *authService) generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *authService) comparePasswordAndHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
