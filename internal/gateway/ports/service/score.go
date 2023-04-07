package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type ScoreService interface {
	GetCVERisk(data domain.ScoreData) (riskValue uint64, scale string, err error)
	GetIdealAge(data domain.ScoreData) (agesRange, scale string, err error)
	ResolveScale(riskValue float64, age int) string
}
