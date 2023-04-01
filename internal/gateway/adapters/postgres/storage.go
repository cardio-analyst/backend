package postgres

import (
	"context"
	"github.com/cardio-analyst/backend/internal/gateway/config"
	storage2 "github.com/cardio-analyst/backend/internal/gateway/ports/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// check whether postgresStorage structure implements the storage.Storage interface
var _ storage2.Storage = (*postgresStorage)(nil)

// postgresStorage implements storage.Storage interface.
type postgresStorage struct {
	conn *pgxpool.Pool

	userRepository            storage2.UserRepository
	sessionRepository         storage2.SessionRepository
	diseasesRepository        storage2.DiseasesRepository
	analysisRepository        storage2.AnalysisRepository
	lifestyleRepository       storage2.LifestyleRepository
	basicIndicatorsRepository storage2.BasicIndicatorsRepository
	scoreRepository           storage2.ScoreRepository
}

func NewStorage(cfg config.PostgresConfig) (*postgresStorage, error) {
	pool, err := pgxpool.Connect(context.Background(), cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &postgresStorage{conn: pool}, nil
}

func (s *postgresStorage) Users() storage2.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = NewUserRepository(s)

	return s.userRepository
}

func (s *postgresStorage) Sessions() storage2.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = NewSessionRepository(s)

	return s.sessionRepository
}

func (s *postgresStorage) Diseases() storage2.DiseasesRepository {
	if s.diseasesRepository != nil {
		return s.diseasesRepository
	}

	s.diseasesRepository = NewDiseasesRepository(s)

	return s.diseasesRepository
}

func (s *postgresStorage) Analyses() storage2.AnalysisRepository {
	if s.analysisRepository != nil {
		return s.analysisRepository
	}

	s.analysisRepository = NewAnalysisRepository(s)

	return s.analysisRepository
}

func (s *postgresStorage) Lifestyles() storage2.LifestyleRepository {
	if s.lifestyleRepository != nil {
		return s.lifestyleRepository
	}

	s.lifestyleRepository = NewLifestyleRepository(s)

	return s.lifestyleRepository
}

func (s *postgresStorage) BasicIndicators() storage2.BasicIndicatorsRepository {
	if s.basicIndicatorsRepository != nil {
		return s.basicIndicatorsRepository
	}

	s.basicIndicatorsRepository = NewBasicIndicatorsRepository(s)

	return s.basicIndicatorsRepository
}

func (s *postgresStorage) Score() storage2.ScoreRepository {
	if s.scoreRepository != nil {
		return s.scoreRepository
	}

	s.scoreRepository = NewScoreRepository(s)

	return s.scoreRepository
}

func (s *postgresStorage) Close() error {
	s.conn.Close()
	return nil
}
