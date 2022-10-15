package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

// UserStorage encapsulates the logic of manipulations on the entity "User" in the database.
type UserStorage interface {
	// SaveUser is a symbiosis of update and insert methods (upsert).
	//
	// If the user is in the database, then the data of the existing user is updated, otherwise the data of the new user is inserted.
	SaveUser(userData models.User) (err error)
	// GetUserByCriteria searches for user in the database according to the criteria fixed in the models.UserCriteria.
	//
	// By the time the method is used, it is assumed that the user definitely exists in the database, so if it is not found,
	// then the method returns an error.
	GetUserByCriteria(criteria models.UserCriteria) (user *models.User, err error)
	// FindUserByCriteria searches for users in the database according to the criteria fixed in the models.UserCriteria.
	//
	// If users with the corresponding criteria are not found, the method returns nil.
	FindUserByCriteria(criteria models.UserCriteria) (user []*models.User, err error)
}
