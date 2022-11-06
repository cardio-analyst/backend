package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

// UserRepository encapsulates the logic of manipulations on the entity "User" in the database.
type UserRepository interface {
	// Save is a symbiosis of update and insert methods (upsert).
	//
	// If the user is in the database, then the data of the existing user is updated, otherwise the data of the new user is inserted.
	Save(userData models.User) (err error)
	// GetByCriteria searches for user in the database according to the criteria fixed in the models.UserCriteria.
	//
	// By the time the method is used, it is assumed that the user definitely exists in the database, so if it is not found,
	// then the method returns an error.
	GetByCriteria(criteria models.UserCriteria) (userData *models.User, err error)
	// FindByCriteria searches for users in the database according to the criteria fixed in the models.UserCriteria.
	//
	// If users with the corresponding criteria are not found, the method returns nil.
	FindByCriteria(criteria models.UserCriteria) (userDataList []*models.User, err error)
}
