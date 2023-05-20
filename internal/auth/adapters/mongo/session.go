package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	domain "github.com/cardio-analyst/backend/internal/auth/domain/model"
)

// SessionRepository implements storage.SessionRepository interface.
type SessionRepository struct {
	storage *Storage
}

func NewSessionRepository(storage *Storage) *SessionRepository {
	return &SessionRepository{
		storage: storage,
	}
}

func (r *SessionRepository) Save(ctx context.Context, session domain.Session) error {
	filter := bson.M{"id": session.UserID}
	update := bson.D{
		{
			Key:   "$set",
			Value: session,
		},
	}
	opts := options.Update().SetUpsert(true)

	result, err := r.storage.sessions.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount != 0 {
		log.Debug("matched and replaced an existing session")
	} else if result.UpsertedCount != 0 {
		log.Debugf("inserted a new session with user ID %v (mongo ID %v)", session.UserID, result.UpsertedID)
	}

	return nil
}

func (r *SessionRepository) GetOne(ctx context.Context, userID uint64) (domain.Session, error) {
	filter := bson.D{
		{
			Key:   "user_id",
			Value: userID,
		},
	}

	var session domain.Session
	if err := r.storage.sessions.FindOne(ctx, filter).Decode(&session); err != nil {
		if err == mongo.ErrNoDocuments {
			return session, domain.ErrSessionNotFound
		}
		return session, err
	}

	return session, nil
}

func (r *SessionRepository) FindOne(ctx context.Context, userID uint64) (*domain.Session, error) {
	filter := bson.D{
		{
			Key:   "user_id",
			Value: userID,
		},
	}

	var session domain.Session
	if err := r.storage.sessions.FindOne(ctx, filter).Decode(&session); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}
