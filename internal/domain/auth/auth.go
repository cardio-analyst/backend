package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	if err := user.Validate(true); err != nil {
		return fmt.Errorf("%w: %v", serviceErrors.ErrInvalidUserData, err)
	}

	criteria := models.UserCriteria{
		Login:             &user.Login,
		Email:             &user.Email,
		CriteriaSeparator: models.CriteriaSeparatorOR,
	}

	users, err := s.db.FindByCriteria(criteria)
	if err != nil {
		return err
	}

	// check if there are no users with a username or email that the actor wants to occupy
	if len(users) > 0 {
		for _, u := range users {
			if u.Login == user.Login {
				return serviceErrors.ErrUserLoginAlreadyOccupied
			} else {
				return serviceErrors.ErrUserEmailAlreadyOccupied
			}
		}
	}

	passwordHash, err := s.generateHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	return s.db.Save(user)
}

func (s *authService) GetToken(credentials models.UserCredentials) (string, error) {
	if err := credentials.Validate(); err != nil {
		return "", fmt.Errorf("%w: %v", serviceErrors.ErrInvalidUserCredentials, err)
	}

	// since we do not know what exactly we are dealing with (login or email), we are looking for two fields
	criteria := models.UserCriteria{
		Login:             &credentials.LoginOrEmail,
		Email:             &credentials.LoginOrEmail,
		CriteriaSeparator: models.CriteriaSeparatorOR,
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

	return s.generateJWTToken(user.ID)
}

func (s *authService) ValidateToken(token string) (uint64, error) {
	parsed, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", serviceErrors.ErrWrongToken, token.Header["alg"])
		}
		return []byte(s.cfg.SigningKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, serviceErrors.ErrWrongToken
		}
		if strings.Contains(err.Error(), "token is expired by") {
			return 0, serviceErrors.ErrTokenIsExpired
		}
		return 0, err
	}

	claims, ok := parsed.Claims.(*tokenClaims)
	if !(parsed.Valid && ok) {
		return 0, serviceErrors.ErrWrongToken
	}

	if time.Now().After(time.Unix(claims.ExpiresAt.Unix(), 0)) {
		return 0, serviceErrors.ErrTokenIsExpired
	}

	return claims.UserID, nil
}

type tokenClaims struct {
	UserID uint64 `json:"userID"`
	jwt.RegisteredClaims
}

func (s *authService) generateJWTToken(userID uint64) (string, error) {
	claims := tokenClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.TokenTTLMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SigningKey))
}

func (s *authService) generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *authService) comparePasswordAndHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
