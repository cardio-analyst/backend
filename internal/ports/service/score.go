package service

import "github.com/cardio-analyst/backend/internal/domain/models"

type ScoreService interface {
	GetCVERisk(data models.ScoreData) (riskValue uint64, scale string, err error)
	GetIdealAge(data models.ScoreData) (agesRange, scale string, err error)
}
