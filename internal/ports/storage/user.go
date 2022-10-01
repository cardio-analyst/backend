package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

// UserStorage TODO
type UserStorage interface {
	// Create TODO
	Create(userData models.User) (err error)
	// GetOneByCriteria TODO (get вместо nil ошибку даёт)
	GetOneByCriteria(criteria models.UserCriteria) (user *models.User, err error)
	// FindOneByCriteria TODO (find может вернуть nil)
	FindOneByCriteria(criteria models.UserCriteria) (user *models.User, err error)
}
