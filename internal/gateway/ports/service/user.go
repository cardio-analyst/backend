package service

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// UserService TODO
type UserService interface {
	// Get TODO
	Get(criteria models.UserCriteria) (userData *models.User, err error)
	// Update TODO
	Update(userData models.User) (err error)
}
