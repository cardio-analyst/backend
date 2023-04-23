package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

// Services represents a layer of business logic.
type Services interface {
	User() UserService
	Auth() AuthService
	Diseases() DiseasesService
	Analysis() AnalysisService
	Lifestyle() LifestyleService
	Questionnaire() QuestionnaireService
	BasicIndicators() BasicIndicatorsService
	Score() ScoreService
	Recommendations() RecommendationsService
	Email() EmailService
	Report(reportType domain.ReportType) ReportService
}
