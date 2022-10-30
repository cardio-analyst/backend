package storage

// Storage represents the data storage layer (database).
type Storage interface {
	// Users allows you to access information about users.
	Users() UserRepository
	// Sessions allows you to access information about sessions.
	Sessions() SessionRepository
	// Diseases allows you to access information about user diseases.
	Diseases() DiseasesRepository
	// Analyses allows you to access information about user analyses.
	Analyses() AnalysisRepository
}
