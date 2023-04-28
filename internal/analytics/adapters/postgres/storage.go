package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
)

type Storage struct {
	conn *pgxpool.Pool

	feedbackRepository storage.FeedbackRepository
}

func NewStorage(dsn string) (*Storage, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &Storage{conn: pool}, nil
}

func (s *Storage) Feedback() storage.FeedbackRepository {
	if s.feedbackRepository != nil {
		return s.feedbackRepository
	}

	s.feedbackRepository = NewFeedbackRepository(s)

	return s.feedbackRepository
}

func (s *Storage) Close() error {
	s.conn.Close()
	return nil
}
