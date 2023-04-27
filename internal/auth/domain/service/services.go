package service

import (
	"github.com/cardio-analyst/backend/internal/auth/config"
	"github.com/cardio-analyst/backend/internal/auth/port/service"
	"github.com/cardio-analyst/backend/internal/auth/port/storage"
)

type Services struct {
	cfg     config.AuthConfig
	storage storage.Storage

	authService       service.AuthService
	userService       service.UserService
	validationService service.ValidationService
}

func NewServices(cfg config.AuthConfig, storage storage.Storage) *Services {
	return &Services{
		cfg:     cfg,
		storage: storage,
	}
}

func (s *Services) Auth() service.AuthService {
	if s.authService != nil {
		return s.authService
	}

	opts := AuthServiceOptions{
		Users:                  s.storage.Users(),
		Sessions:               s.storage.Sessions(),
		AccessTokenTTLSeconds:  s.cfg.AccessToken.TokenTTLSec,
		RefreshTokenTTLSeconds: s.cfg.RefreshToken.TokenTTLSec,
		AccessTokenSigningKey:  s.cfg.AccessToken.SigningKey,
		RefreshTokenSigningKey: s.cfg.RefreshToken.SigningKey,
		SecretKeySigningKey:    s.cfg.SecretKeySigningKey,
	}

	s.authService = NewAuthService(opts)

	return s.authService
}

func (s *Services) User() service.UserService {
	if s.userService != nil {
		return s.userService
	}

	s.userService = NewUserService(s.storage.Users())

	return s.userService
}

func (s *Services) Validation() service.ValidationService {
	if s.validationService != nil {
		return s.validationService
	}

	s.validationService = NewValidationService()

	return s.validationService
}
