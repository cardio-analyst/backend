package service

import "github.com/cardio-analyst/backend/internal/domain/models"

// UserService TODO
type UserService interface {
	// Get TODO
	Get(criteria models.UserCriteria) (user *models.User, err error)
	// Update TODO
	Update(user models.User) (err error)
}
