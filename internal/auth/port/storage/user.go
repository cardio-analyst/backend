package storage

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

// UserRepository encapsulates the logic of manipulations on the entity "User" in the database.
type UserRepository interface {
	// Save is a symbiosis of update and insert methods (upsert).
	//
	// If the user is in the database, then the data of the existing user is updated, otherwise the data of the new user is inserted.
	Save(ctx context.Context, user model.User) (err error)
	// GetOneByCriteria searches for user in the database according to the criteria fixed in the model.UserCriteria.
	//
	// By the time the method is used, it is assumed that the user definitely exists in the database, so if it is not found,
	// then the method returns an error.
	GetOneByCriteria(ctx context.Context, criteria model.UserCriteria) (user model.User, err error)
	// FindAllByCriteria searches for users in the database according to the criteria fixed in the model.UserCriteria.
	//
	// If users with the corresponding criteria are not found, the method returns nil.
	FindAllByCriteria(ctx context.Context, criteria model.UserCriteria) (users []model.User, err error)
}
