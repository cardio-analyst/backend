package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/config"
	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether services structure implements the service.Services interface
var _ service.Services = (*Services)(nil)

// Services implements service.Services interface.
type Services struct {
	cfg     config.Config
	storage storage.Storage

	emailPublisher    client.EmailPublisher
	feedbackPublisher client.FeedbackPublisher

	authClient      client.Auth
	analyticsClient client.Analytics

	userService            service.UserService
	authService            service.AuthService
	diseasesService        service.DiseasesService
	analysisService        service.AnalysisService
	lifestyleService       service.LifestyleService
	questionnaireService   service.QuestionnaireService
	basicIndicatorsService service.BasicIndicatorsService
	scoreService           service.ScoreService
	recommendationsService service.RecommendationsService
	emailService           service.EmailService
	feedbackService        service.FeedbackService
	reportService          service.ReportService
}

type ServicesOptions struct {
	Config            config.Config
	Storage           storage.Storage
	EmailPublisher    client.EmailPublisher
	FeedbackPublisher client.FeedbackPublisher
	AuthClient        client.Auth
	AnalyticsClient   client.Analytics
}

func NewServices(opts ServicesOptions) *Services {
	return &Services{
		cfg:               opts.Config,
		storage:           opts.Storage,
		emailPublisher:    opts.EmailPublisher,
		feedbackPublisher: opts.FeedbackPublisher,
		authClient:        opts.AuthClient,
		analyticsClient:   opts.AnalyticsClient,
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

func (s *Services) Questionnaire() service.QuestionnaireService {
	if s.questionnaireService != nil {
		return s.questionnaireService
	}

	s.questionnaireService = NewQuestionnaireService(s.storage.Questionnaire())

	return s.questionnaireService
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

	s.emailService = NewEmailService(s.emailPublisher)

	return s.emailService
}

func (s *Services) Feedback() service.FeedbackService {
	if s.feedbackService != nil {
		return s.feedbackService
	}

	s.feedbackService = NewFeedbackService(s.feedbackPublisher, s.analyticsClient)

	return s.feedbackService
}

func (s *Services) Report() service.ReportService {
	if s.reportService != nil {
		return s.reportService
	}

	s.reportService = NewPDFReportService(
		s.Recommendations(),
		s.storage.Analyses(),
		s.storage.BasicIndicators(),
		s.authClient,
	)

	return s.reportService
}
