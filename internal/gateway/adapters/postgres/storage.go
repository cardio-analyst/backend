package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether Storage structure implements the storage.Storage interface
var _ storage.Storage = (*Storage)(nil)

// Storage implements storage.Storage interface.
type Storage struct {
	conn *pgxpool.Pool

	diseasesRepository        storage.DiseasesRepository
	analysisRepository        storage.AnalysisRepository
	lifestyleRepository       storage.LifestyleRepository
	questionnaireRepository   storage.QuestionnaireRepository
	basicIndicatorsRepository storage.BasicIndicatorsRepository
	scoreRepository           storage.ScoreRepository
}

func NewStorage(dsn string) (*Storage, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &Storage{conn: pool}, nil
}

func (s *Storage) Diseases() storage.DiseasesRepository {
	if s.diseasesRepository != nil {
		return s.diseasesRepository
	}

	s.diseasesRepository = NewDiseasesRepository(s)

	return s.diseasesRepository
}

func (s *Storage) Analyses() storage.AnalysisRepository {
	if s.analysisRepository != nil {
		return s.analysisRepository
	}

	s.analysisRepository = NewAnalysisRepository(s)

	return s.analysisRepository
}

func (s *Storage) Lifestyles() storage.LifestyleRepository {
	if s.lifestyleRepository != nil {
		return s.lifestyleRepository
	}

	s.lifestyleRepository = NewLifestyleRepository(s)

	return s.lifestyleRepository
}

func (s *Storage) Questionnaire() storage.QuestionnaireRepository {
	if s.questionnaireRepository != nil {
		return s.questionnaireRepository
	}

	s.questionnaireRepository = NewQuestionnaireRepository(s)

	return s.questionnaireRepository
}

func (s *Storage) BasicIndicators() storage.BasicIndicatorsRepository {
	if s.basicIndicatorsRepository != nil {
		return s.basicIndicatorsRepository
	}

	s.basicIndicatorsRepository = NewBasicIndicatorsRepository(s)

	return s.basicIndicatorsRepository
}

func (s *Storage) Score() storage.ScoreRepository {
	if s.scoreRepository != nil {
		return s.scoreRepository
	}

	s.scoreRepository = NewScoreRepository(s)

	return s.scoreRepository
}

func (s *Storage) Close() error {
	s.conn.Close()
	return nil
}
