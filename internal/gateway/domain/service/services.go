package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/config"
	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether services structure implements the service.Services interface
var _ service.Services = (*Services)(nil)

// Services implements service.Services interface.
type Services struct {
	cfg            config.Config
	storage        storage.Storage
	rabbitMQClient client.RabbitMQPublisher
	authClient     client.Auth

	userService            service.UserService
	authService            service.AuthService
	diseasesService        service.DiseasesService
	analysisService        service.AnalysisService
	lifestyleService       service.LifestyleService
	basicIndicatorsService service.BasicIndicatorsService
	scoreService           service.ScoreService
	recommendationsService service.RecommendationsService
	emailService           service.EmailService

	reportServices ReportServices
}

type ReportServices struct {
	PDF service.ReportService
}

func NewServices(
	cfg config.Config,
	storage storage.Storage,
	rabbitMQClient client.RabbitMQPublisher,
	authClient client.Auth,
) *Services {
	return &Services{
		cfg:            cfg,
		storage:        storage,
		rabbitMQClient: rabbitMQClient,
		authClient:     authClient,
	}
}

func (s *Services) User() service.UserService {
	if s.userService != nil {
		return s.userService
	}

	s.userService = NewUserService(s.authClient)

	return s.userService
}

func (s *Services) Auth() service.AuthService {
	if s.authService != nil {
		return s.authService
	}

	s.authService = NewAuthService(s.authClient)

	return s.authService
}

func (s *Services) Diseases() service.DiseasesService {
	if s.diseasesService != nil {
		return s.diseasesService
	}

	s.diseasesService = NewDiseasesService(s.storage.Diseases())

	return s.diseasesService
}

func (s *Services) Analysis() service.AnalysisService {
	if s.analysisService != nil {
		return s.analysisService
	}

	s.analysisService = NewAnalysisService(s.storage.Analyses())

	return s.analysisService
}

func (s *Services) Lifestyle() service.LifestyleService {
	if s.lifestyleService != nil {
		return s.lifestyleService
	}

	s.lifestyleService = NewLifestyleService(s.storage.Lifestyles())

	return s.lifestyleService
}

func (s *Services) BasicIndicators() service.BasicIndicatorsService {
	if s.basicIndicatorsService != nil {
		return s.basicIndicatorsService
	}

	s.basicIndicatorsService = NewBasicIndicatorsService(s.storage.BasicIndicators())

	return s.basicIndicatorsService
}

func (s *Services) Score() service.ScoreService {
	if s.scoreService != nil {
		return s.scoreService
	}

	s.scoreService = NewScoreService(s.storage.Score())

	return s.scoreService
}

func (s *Services) Recommendations() service.RecommendationsService {
	if s.recommendationsService != nil {
		return s.recommendationsService
	}

	s.recommendationsService = NewRecommendationsService(
		s.cfg.Recommendations,
		s.storage.Diseases(),
		s.storage.BasicIndicators(),
		s.storage.Lifestyles(),
		s.storage.Score(),
		s.authClient,
	)

	return s.recommendationsService
}

func (s *Services) Email() service.EmailService {
	if s.emailService != nil {
		return s.emailService
	}

	s.emailService = NewEmailService(s.rabbitMQClient)

	return s.emailService
}

func (s *Services) Report(reportType domain.ReportType) service.ReportService {
	switch reportType {
	case domain.PDF:
		if s.reportServices.PDF != nil {
			return s.reportServices.PDF
		}

		s.reportServices.PDF = NewPDFReportService(
			s.Recommendations(),
			s.storage.Analyses(),
			s.storage.BasicIndicators(),
			s.authClient,
		)

		return s.reportServices.PDF
	default:
		return nil
	}
}
