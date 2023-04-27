package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/cardio-analyst/backend/pkg/model"
)

type secretKeyClaims struct {
	UserLogin string `json:"user_login"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateSecretKey(userLogin, userEmail string) (string, error) {
	return s.generateSecretKey(userLogin, userEmail, s.secretKeySigningKey)
}

func (s *AuthService) generateSecretKey(userLogin, userEmail, signingKey string) (string, error) {
	claims := secretKeyClaims{
		UserLogin: userLogin,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	secretKey := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return secretKey.SignedString([]byte(signingKey))
}

func (s *AuthService) VerifySecretKey(user model.User) error {
	if user.Role == model.UserRoleCustomer {
		return nil
	}
	login, email, err := s.validateSecretKey(user.SecretKey, s.secretKeySigningKey)
	if err != nil {
		return err
	}
	if user.Login != login {
		return model.ErrWrongSecretKey
	}
	if user.Email != email {
		return model.ErrWrongSecretKey
	}
	return nil
}

func (s *AuthService) validateSecretKey(secretKey, signingKey string) (string, string, error) {
	parsed, err := jwt.ParseWithClaims(secretKey, &secretKeyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", model.ErrWrongSecretKey, token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "signature is invalid") {
			return "", "", model.ErrWrongSecretKey
		}
		return "", "", err
	}

	claims, ok := parsed.Claims.(*secretKeyClaims)
	if !(parsed.Valid && ok) {
		return "", "", model.ErrWrongSecretKey
	}

	return claims.UserLogin, claims.UserEmail, nil
}
