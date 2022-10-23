package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/cardio-analyst/backend/internal/config"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

// check whether postgresStorage structure implements the storage.Storage interface
var _ storage.Storage = (*postgresStorage)(nil)

// postgresStorage implements storage.Storage interface.
type postgresStorage struct {
	conn *pgxpool.Pool

	userRepository     storage.UserRepository
	sessionRepository  storage.SessionRepository
	diseasesRepository storage.DiseasesRepository
}

func NewStorage(cfg config.PostgresConfig) (*postgresStorage, error) {
	pool, err := pgxpool.Connect(context.Background(), cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &postgresStorage{conn: pool}, nil
}

func (s *postgresStorage) Users() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = NewUserRepository(s)

	return s.userRepository
}

func (s *postgresStorage) Sessions() storage.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = NewSessionRepository(s)

	return s.sessionRepository
}

func (s *postgresStorage) Diseases() storage.DiseasesRepository {
	if s.diseasesRepository != nil {
		return s.diseasesRepository
	}

	s.diseasesRepository = NewDiseasesRepository(s)

	return s.diseasesRepository
}

func (s *postgresStorage) Close() error {
	s.conn.Close()
	return nil
}
