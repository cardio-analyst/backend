package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type BasicIndicatorsService interface {
	Create(basicIndicatorsData domain.BasicIndicators) (err error)
	Update(basicIndicatorsData domain.BasicIndicators) (err error)
	FindAll(userID uint64) (basicIndicatorsDataList []*domain.BasicIndicators, err error)
}
