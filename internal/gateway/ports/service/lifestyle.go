package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type LifestyleService interface {
	Update(lifestyleData domain.Lifestyle) (err error)
	Get(userID uint64) (lifestyle *domain.Lifestyle, err error)
}
