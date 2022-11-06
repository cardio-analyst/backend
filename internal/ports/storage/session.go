package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

// SessionRepository encapsulates the logic of manipulations on the entity "Session" in the database.
type SessionRepository interface {
	// Save is a symbiosis of update and insert methods (upsert).
	//
	// If the session is in the database, then the data of the existing session is updated, otherwise the data of the new session is inserted.
	Save(sessionData models.Session) (err error)
	// Get searches for session in the database according to the user id.
	//
	// By the time the method is used, it is assumed that the session definitely exists in the database, so if it is not found,
	// then the method returns an error.
	Get(userID uint64) (sessionData *models.Session, err error)
	// Find searches for session in the database according to the user id.
	//
	// If session is not found, the method returns nil.
	Find(userID uint64) (sessionData *models.Session, err error)
}
