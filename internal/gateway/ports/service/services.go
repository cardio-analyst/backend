package service

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
	Feedback() FeedbackService
	Report() ReportService
	Statistics() StatisticsService
}
