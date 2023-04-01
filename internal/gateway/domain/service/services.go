package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/config"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	service2 "github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/smtp"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether services structure implements the service.Services interface
var _ service2.Services = (*services)(nil)

// services implements service.Services interface.
type services struct {
	cfg        config.ServicesConfig
	storage    storage.Storage
	smtpClient smtp.Client

	userService            service2.UserService
	authService            service2.AuthService
	diseasesService        service2.DiseasesService
	analysisService        service2.AnalysisService
	lifestyleService       service2.LifestyleService
	basicIndicatorsService service2.BasicIndicatorsService
	scoreService           service2.ScoreService
	recommendationsService service2.RecommendationsService
	emailService           service2.EmailService

	reportServices reportServices
}

type reportServices struct {
	PDF service2.ReportService
}

func NewServices(
	cfg config.ServicesConfig,
	storage storage.Storage,
	smtpClient smtp.Client,
) *services {
	return &services{
		cfg:        cfg,
		storage:    storage,
		smtpClient: smtpClient,
	}
}

func (s *services) User() service2.UserService {
	if s.userService != nil {
		return s.userService
	}

	s.userService = NewUserService(s.storage.Users())

	return s.userService
}

func (s *services) Auth() service2.AuthService {
	if s.authService != nil {
		return s.authService
	}

	s.authService = NewAuthService(s.cfg.Auth, s.storage.Users(), s.storage.Sessions())

	return s.authService
}

func (s *services) Diseases() service2.DiseasesService {
	if s.diseasesService != nil {
		return s.diseasesService
	}

	s.diseasesService = NewDiseasesService(s.storage.Diseases())

	return s.diseasesService
}

func (s *services) Analysis() service2.AnalysisService {
	if s.analysisService != nil {
		return s.analysisService
	}

	s.analysisService = NewAnalysisService(s.storage.Analyses())

	return s.analysisService
}

func (s *services) Lifestyle() service2.LifestyleService {
	if s.lifestyleService != nil {
		return s.lifestyleService
	}

	s.lifestyleService = NewLifestyleService(s.storage.Lifestyles())

	return s.lifestyleService
}

func (s *services) BasicIndicators() service2.BasicIndicatorsService {
	if s.basicIndicatorsService != nil {
		return s.basicIndicatorsService
	}

	s.basicIndicatorsService = NewBasicIndicatorsService(s.storage.BasicIndicators())

	return s.basicIndicatorsService
}

func (s *services) Score() service2.ScoreService {
	if s.scoreService != nil {
		return s.scoreService
	}

	s.scoreService = NewScoreService(s.storage.Score())

	return s.scoreService
}

func (s *services) Recommendations() service2.RecommendationsService {
	if s.recommendationsService != nil {
		return s.recommendationsService
	}

	s.recommendationsService = NewRecommendationsService(
		s.cfg.Recommendations,
		s.storage.Diseases(),
		s.storage.BasicIndicators(),
		s.storage.Lifestyles(),
		s.storage.Score(),
		s.storage.Users(),
	)

	return s.recommendationsService
}

func (s *services) Email() service2.EmailService {
	if s.emailService != nil {
		return s.emailService
	}

	s.emailService = NewEmailService(s.smtpClient)

	return s.emailService
}

func (s *services) Report(reportType models.ReportType) service2.ReportService {
	switch reportType {
	case models.PDF:
		if s.reportServices.PDF != nil {
			return s.reportServices.PDF
		}

		s.reportServices.PDF = NewPDFReportService(
			s.Recommendations(),
			s.storage.Analyses(),
			s.storage.BasicIndicators(),
			s.storage.Users(),
		)

		return s.reportServices.PDF
	default:
		return nil
	}
}
