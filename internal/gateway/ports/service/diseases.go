package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type DiseasesService interface {
	Update(diseasesData domain.Diseases) (err error)
	Get(userID uint64) (diseasesData *domain.Diseases, err error)
}
