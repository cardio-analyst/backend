package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cardio-analyst/backend/internal/auth/port/storage"
)

const (
	collectionCounters = "counters"
	collectionUsers    = "users"
	collectionSessions = "sessions"
)

type Storage struct {
	counters *mongo.Collection
	users    *mongo.Collection
	sessions *mongo.Collection

	userRepository    storage.UserRepository
	sessionRepository storage.SessionRepository
}

func NewStorage(dsn, dbname string) (*Storage, error) {
	clientOptions := options.Client().ApplyURI(dsn)

	mongoClient, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if err = mongoClient.Connect(ctx); err != nil {
		return nil, err
	}

	counters := mongoClient.Database(dbname).Collection(collectionCounters)
	if err = initCounter(ctx, counters, counterNameUserID); err != nil {
		return nil, err
	}

	return &Storage{
		counters: counters,
		users:    mongoClient.Database(dbname).Collection(collectionUsers),
		sessions: mongoClient.Database(dbname).Collection(collectionSessions),
	}, nil
}

func (s *Storage) Users() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = NewUserRepository(s)

	return s.userRepository
}

func (s *Storage) Sessions() storage.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = NewSessionRepository(s)

	return s.sessionRepository
}

func (s *Storage) Close() error {
	return s.counters.Database().Client().Disconnect(context.Background())
}
