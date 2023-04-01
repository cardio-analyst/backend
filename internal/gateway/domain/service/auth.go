package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cardio-analyst/backend/internal/gateway/config"
	errors2 "github.com/cardio-analyst/backend/internal/gateway/domain/errors"
	models2 "github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	storage2 "github.com/cardio-analyst/backend/internal/gateway/ports/storage"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"
)

// check whether authService structure implements the service.AuthService interface
var _ service.AuthService = (*authService)(nil)

// authService implements service.AuthService interface.
type authService struct {
	cfg config.AuthConfig

	users    storage2.UserRepository
	sessions storage2.SessionRepository
}

func NewAuthService(
	cfg config.AuthConfig,
	users storage2.UserRepository,
	sessions storage2.SessionRepository,
) *authService {
	return &authService{
		cfg:      cfg,
		users:    users,
		sessions: sessions,
	}
}

func (s *authService) RegisterUser(user models2.User) error {
	if err := user.Validate(true); err != nil {
		return err
	}

	criteria := models2.UserCriteria{
		Login:             &user.Login,
		Email:             &user.Email,
		CriteriaSeparator: models2.CriteriaSeparatorOR,
	}

	users, err := s.users.FindByCriteria(criteria)
	if err != nil {
		return err
	}

	// check if there are no users with a username or email that the actor wants to occupy
	if len(users) > 0 {
		for _, u := range users {
			if u.Login == user.Login {
				return errors2.ErrUserLoginAlreadyOccupied
			} else {
				return errors2.ErrUserEmailAlreadyOccupied
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

func (s *authService) GetTokens(credentials models2.UserCredentials, userIP string) (*models2.Tokens, error) {
	if err := credentials.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errors2.ErrInvalidUserCredentials, err)
	}

	// since we do not know what exactly we are dealing with (login or email), we are looking for two fields
	criteria := models2.UserCriteria{
		Login:             &credentials.LoginOrEmail,
		Email:             &credentials.LoginOrEmail,
		CriteriaSeparator: models2.CriteriaSeparatorOR,
	}

	user, err := s.users.GetByCriteria(criteria)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrWrongCredentials
		}
		return nil, err
	}

	match, err := s.comparePasswordAndHash(credentials.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors2.ErrWrongCredentials
	}

	tokens, err := s.getTokens(user.ID)
	if err != nil {
		return nil, err
	}

	dbSession, err := s.sessions.Find(user.ID)
	if err != nil {
		return nil, err
	}

	session := models2.Session{
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

func (s *authService) RefreshTokens(refreshToken, userIP string) (*models2.Tokens, error) {
	userID, err := s.validateToken(refreshToken, s.cfg.RefreshToken.SigningKey)
	if err != nil {
		return nil, err
	}

	session, err := s.sessions.Get(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrSessionNotFound
		}
		return nil, err
	}

	if !session.IsIPAllowed(userIP) {
		return nil, errors2.ErrIPIsNotInWhitelist
	}

	if refreshToken != session.RefreshToken {
		return nil, errors2.ErrWrongToken
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

func (s *authService) getTokens(userID uint64) (*models2.Tokens, error) {
	accessToken, err := s.generateJWTToken(userID, time.Duration(s.cfg.AccessToken.TokenTTLSec), s.cfg.AccessToken.SigningKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateJWTToken(userID, time.Duration(s.cfg.RefreshToken.TokenTTLSec), s.cfg.RefreshToken.SigningKey)
	if err != nil {
		return nil, err
	}

	return &models2.Tokens{
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
			return nil, fmt.Errorf("%w: %v", errors2.ErrWrongToken, token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "signature is invalid") {
			return 0, errors2.ErrWrongToken
		}
		if strings.Contains(err.Error(), "token is expired by") {
			return 0, errors2.ErrTokenIsExpired
		}
		return 0, err
	}

	claims, ok := parsed.Claims.(*tokenClaims)
	if !(parsed.Valid && ok) {
		return 0, errors2.ErrWrongToken
	}

	if time.Now().After(time.Unix(claims.ExpiresAt.Unix(), 0)) {
		return 0, errors2.ErrTokenIsExpired
	}

	return claims.UserID, nil
}

func (s *authService) generateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *authService) comparePasswordAndHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
