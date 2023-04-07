package storage

// Storage represents the data storage layer (database).
type Storage interface {
	Diseases() DiseasesRepository
	Analyses() AnalysisRepository
	Lifestyles() LifestyleRepository
	BasicIndicators() BasicIndicatorsRepository
	Score() ScoreRepository
}
