package service

// Services represents a layer of business logic.
type Services interface {
	// User contains the methods of business logic of working with users.
	User() UserService
	// Auth contains the methods of business logic of working with authorization.
	Auth() AuthService
	// Diseases contains the methods of business logic of working with user diseases.
	Diseases() DiseasesService
	// Analysis contains the methods of business logic of working with user analyses.
	Analysis() AnalysisService
	// Lifestyle TODO
	Lifestyle() LifestyleService
	// BasicIndicators TODO
	BasicIndicators() BasicIndicatorsService
	// Score TODO
	Score() ScoreService
}
