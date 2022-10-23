package service

import (
	"github.com/cardio-analyst/backend/internal/config"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// check whether services structure implements the service.Services interface
var _ service.Services = (*services)(nil)

// services implements service.Services interface.
type services struct {
	cfg     config.ServicesConfig
	storage storage.Storage

	userService    service.UserService
	authService    service.AuthService
	diseaseService service.DiseaseService
}

func NewServices(cfg config.ServicesConfig, storage storage.Storage) *services {
	return &services{
		cfg:     cfg,
		storage: storage,
	}
}

func (s *services) User() service.UserService {
	if s.userService != nil {
		return s.userService
	}

	s.userService = NewUserService(s.storage.Users())

	return s.userService
}

func (s *services) Auth() service.AuthService {
	if s.authService != nil {
		return s.authService
	}

	s.authService = NewAuthService(s.cfg.Auth, s.storage.Users(), s.storage.Sessions())

	return s.authService
}

func (s *services) Disease() service.DiseaseService {
	if s.diseaseService != nil {
		return s.diseaseService
	}

	s.diseaseService = NewDiseaseService(s.storage.Diseases())

	return s.diseaseService
}
