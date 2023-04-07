package storage

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type LifestyleRepository interface {
	Update(lifestyleData domain.Lifestyle) (err error)
	Get(userID uint64) (lifestyleData *domain.Lifestyle, err error)
}
