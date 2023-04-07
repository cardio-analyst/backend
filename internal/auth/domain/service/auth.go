package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"

	domain "github.com/cardio-analyst/backend/internal/auth/domain/model"
	"github.com/cardio-analyst/backend/internal/auth/port/storage"
	"github.com/cardio-analyst/backend/pkg/model"
)

type AuthService struct {
	users    storage.UserRepository
	sessions storage.SessionRepository

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	accessTokenSigningKey  string
	refreshTokenSigningKey string
}

type AuthServiceOptions struct {
	Users    storage.UserRepository
	Sessions storage.SessionRepository

	AccessTokenTTLSeconds  int64
	RefreshTokenTTLSeconds int64

	AccessTokenSigningKey  string
	RefreshTokenSigningKey string
}

func NewAuthService(opts AuthServiceOptions) *AuthService {
	return &AuthService{
		users:                  opts.Users,
		sessions:               opts.Sessions,
		accessTokenTTL:         time.Duration(opts.AccessTokenTTLSeconds) * time.Second,
		refreshTokenTTL:        time.Duration(opts.RefreshTokenTTLSeconds) * time.Second,
		accessTokenSigningKey:  opts.AccessTokenSigningKey,
		refreshTokenSigningKey: opts.RefreshTokenSigningKey,
	}
}

func (s *AuthService) GetTokens(ctx context.Context, credentials model.Credentials, userIP string) (model.Tokens, error) {
	// since we do not know what exactly we are dealing with (login or email), we are looking for two fields
	criteria := model.UserCriteria{
		Login:             credentials.Login,
		Email:             credentials.Email,
		CriteriaSeparator: model.CriteriaSeparatorOR,
	}

	user, err := s.users.GetOneByCriteria(ctx, criteria)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return model.Tokens{}, model.ErrWrongCredentials
		}
		return model.Tokens{}, err
	}

	match, err := comparePasswordAndHash(credentials.Password, user.Password)
	if err != nil {
		return model.Tokens{}, err
	}
	if !match {
		return model.Tokens{}, model.ErrWrongCredentials
	}

	tokens, err := s.getTokens(user.ID, user.Role)
	if err != nil {
		return model.Tokens{}, err
	}

	dbSession, err := s.sessions.FindOne(ctx, user.ID)
	if err != nil {
		return model.Tokens{}, err
	}

	session := domain.Session{
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

	if err = s.sessions.Save(ctx, session); err != nil {
		return model.Tokens{}, err
	}

	return tokens, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken, userIP string) (model.Tokens, error) {
	userID, userRole, err := s.validateToken(refreshToken, s.refreshTokenSigningKey)
	if err != nil {
		return model.Tokens{}, err
	}

	session, err := s.sessions.GetOne(ctx, userID)
	if err != nil {
		return model.Tokens{}, err
	}

	if !session.IsIPAllowed(userIP) {
		return model.Tokens{}, model.ErrIPIsNotInWhitelist
	}

	if refreshToken != session.RefreshToken {
		return model.Tokens{}, model.ErrWrongToken
	}

	newTokens, err := s.getTokens(userID, userRole)
	if err != nil {
		return model.Tokens{}, err
	}

	session.RefreshToken = newTokens.RefreshToken

	if err = s.sessions.Save(ctx, session); err != nil {
		return model.Tokens{}, err
	}

	return newTokens, nil
}

func (s *AuthService) IdentifyUser(_ context.Context, token string) (uint64, model.UserRole, error) {
	return s.validateToken(token, s.accessTokenSigningKey)
}

func comparePasswordAndHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func (s *AuthService) getTokens(userID uint64, userRole model.UserRole) (model.Tokens, error) {
	accessToken, err := s.generateJWTToken(userID, userRole, s.accessTokenTTL, s.accessTokenSigningKey)
	if err != nil {
		return model.Tokens{}, err
	}

	refreshToken, err := s.generateJWTToken(userID, userRole, s.refreshTokenTTL, s.refreshTokenSigningKey)
	if err != nil {
		return model.Tokens{}, err
	}

	return model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type tokenClaims struct {
	UserID   uint64         `json:"user_id"`
	UserRole model.UserRole `json:"user_role"`
	jwt.RegisteredClaims
}

func (s *AuthService) generateJWTToken(userID uint64, userRole model.UserRole, ttl time.Duration, signingKey string) (string, error) {
	claims := tokenClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) validateToken(token, signingKey string) (uint64, model.UserRole, error) {
	parsed, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", model.ErrWrongToken, token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "signature is invalid") {
			return 0, "", model.ErrWrongToken
		}
		if strings.Contains(err.Error(), "token is expired by") {
			return 0, "", model.ErrTokenIsExpired
		}
		return 0, "", err
	}

	claims, ok := parsed.Claims.(*tokenClaims)
	if !(parsed.Valid && ok) {
		return 0, "", model.ErrWrongToken
	}

	if time.Now().After(time.Unix(claims.ExpiresAt.Unix(), 0)) {
		return 0, "", model.ErrTokenIsExpired
	}

	return claims.UserID, claims.UserRole, nil
}
