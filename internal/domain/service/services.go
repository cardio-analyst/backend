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

	userService            service.UserService
	authService            service.AuthService
	diseasesService        service.DiseasesService
	analysisService        service.AnalysisService
	lifestyleService       service.LifestyleService
	basicIndicatorsService service.BasicIndicatorsService
	scoreService           service.ScoreService
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

func (s *services) Diseases() service.DiseasesService {
	if s.diseasesService != nil {
		return s.diseasesService
	}

	s.diseasesService = NewDiseasesService(s.storage.Diseases())

	return s.diseasesService
}

func (s *services) Analysis() service.AnalysisService {
	if s.analysisService != nil {
		return s.analysisService
	}

	s.analysisService = NewAnalysisService(s.storage.Analyses())

	return s.analysisService
}

func (s *services) Lifestyle() service.LifestyleService {
	if s.lifestyleService != nil {
		return s.lifestyleService
	}

	s.lifestyleService = NewLifestyleService(s.storage.Lifestyles())

	return s.lifestyleService
}

func (s *services) BasicIndicators() service.BasicIndicatorsService {
	if s.basicIndicatorsService != nil {
		return s.basicIndicatorsService
	}

	s.basicIndicatorsService = NewBasicIndicatorsService(s.storage.BasicIndicators())

	return s.basicIndicatorsService
}

func (s *services) Score() service.ScoreService {
	if s.scoreService != nil {
		return s.scoreService
	}

	s.scoreService = NewScoreService(s.storage.Score())

	return s.scoreService
}
