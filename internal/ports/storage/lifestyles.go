package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

type LifestyleRepository interface {
	Update(lifestyleData models.Lifestyle) (err error)

	Get(userID uint64) (lifestyle *models.Lifestyle, err error)
}

