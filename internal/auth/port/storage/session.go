package storage

import (
	"context"

	domain "github.com/cardio-analyst/backend/internal/auth/domain/model"
)

// SessionRepository encapsulates the logic of manipulations on the entity "Session" in the database.
type SessionRepository interface {
	// Save is a symbiosis of update and insert methods (upsert).
	//
	// If the session is in the database, then the data of the existing session is updated, otherwise the data of the new session is inserted.
	Save(ctx context.Context, session domain.Session) (err error)
	// GetOne searches for session in the database according to the user id.
	//
	// By the time the method is used, it is assumed that the session definitely exists in the database, so if it is not found,
	// then the method returns an error.
	GetOne(ctx context.Context, userID uint64) (session domain.Session, err error)
	// FindOne searches for session in the database according to the user id.
	//
	// If session is not found, the method returns nil.
	FindOne(ctx context.Context, userID uint64) (session *domain.Session, err error)
}
