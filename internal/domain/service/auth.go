package service

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

// check whether authService structure implements the service.AuthService interface
var _ service.AuthService = (*authService)(nil)

// authService implements service.AuthService interface.
type authService struct {
	cfg config.AuthConfig

	users    storage.UserRepository
	sessions storage.SessionRepository
}

func NewAuthService(
	cfg config.AuthConfig,
	users storage.UserRepository,
	sessions storage.SessionRepository,
) *authService {
	return &authService{
		cfg:      cfg,
		users:    users,
		sessions: sessions,
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

	users, err := s.users.FindByCriteria(criteria)
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

	return s.users.Save(user)
}

func (s *authService) GetTokens(credentials models.UserCredentials, userIP string) (*models.Tokens, error) {
	if err := credentials.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", serviceErrors.ErrInvalidUserCredentials, err)
	}

	// since we do not know what exactly we are dealing with (login or email), we are looking for two fields
	criteria := models.UserCriteria{
		Login:             &credentials.LoginOrEmail,
		Email:             &credentials.LoginOrEmail,
		CriteriaSeparator: models.CriteriaSeparatorOR,
	}

	user, err := s.users.GetByCriteria(criteria)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrWrongCredentials
		}
		return nil, err
	}

	match, err := s.comparePasswordAndHash(credentials.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, serviceErrors.ErrWrongCredentials
	}

	tokens, err := s.getTokens(user.ID)
	if err != nil {
		return nil, err
	}

	dbSession, err := s.sessions.Find(user.ID)
	if err != nil {
		return nil, err
	}

	session := models.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
	}

	if dbSession == nil {
		session.Whitelist = []string{userIP}
	} else {
		if dbSession.IsIPAllowed(userIP) {
			session.Whitelist = dbSession.Whitelist
		} else {
			session.Whitelist = append(dbSession.Whitelist, userIP)
		}
	}

	if err = s.sessions.Save(session); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *authService) RefreshTokens(refreshToken, userIP string) (*models.Tokens, error) {
	userID, err := s.validateToken(refreshToken, s.cfg.RefreshToken.SigningKey)
	if err != nil {
		return nil, err
	}

	session, err := s.sessions.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceErrors.ErrSessionNotFound
		}
		return nil, err
	}

	if !session.IsIPAllowed(userIP) {
		return nil, serviceErrors.ErrIPIsNotInWhitelist
	}

	if refreshToken != session.RefreshToken {
		return nil, serviceErrors.ErrWrongToken
	}

	newTokens, err := s.getTokens(userID)
	if err != nil {
		return nil, err
	}

	session.RefreshToken = newTokens.RefreshToken

	if err = s.sessions.Save(*session); err != nil {
		return nil, err
	}

	return newTokens, nil
}

func (s *authService) getTokens(userID uint64) (*models.Tokens, error) {
	accessToken, err := s.generateJWTToken(userID, time.Duration(s.cfg.AccessToken.TokenTTLSec), s.cfg.AccessToken.SigningKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateJWTToken(userID, time.Duration(s.cfg.RefreshToken.TokenTTLSec), s.cfg.RefreshToken.SigningKey)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type tokenClaims struct {
	UserID uint64 `json:"userID"`
	jwt.RegisteredClaims
}

func (s *authService) generateJWTToken(userID uint64, ttl time.Duration, signingKey string) (string, error) {
	claims := tokenClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func (s *authService) ValidateAccessToken(token string) (uint64, error) {
	return s.validateToken(token, s.cfg.AccessToken.SigningKey)
}

func (s *authService) validateToken(token, signingKey string) (uint64, error) {
	parsed, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", serviceErrors.ErrWrongToken, token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "signature is invalid") {
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

func (s *authService) generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *authService) comparePasswordAndHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
