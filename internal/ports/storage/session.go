package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

// SessionStorage encapsulates the logic of manipulations on the entity "Session" in the database.
type SessionStorage interface {
	// SaveSession is a symbiosis of update and insert methods (upsert).
	//
	// If the session is in the database, then the data of the existing session is updated, otherwise the data of the new session is inserted.
	SaveSession(sessionData models.Session) (err error)
	// GetSession searches for session in the database according to the user id.
	//
	// By the time the method is used, it is assumed that the session definitely exists in the database, so if it is not found,
	// then the method returns an error.
	GetSession(userID uint64) (session *models.Session, err error)
	// FindSession searches for session in the database according to the user id.
	//
	// If session is not found, the method returns nil.
	FindSession(userID uint64) (session *models.Session, err error)
}
