package storage

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

type ScoreRepository interface {
	GetCVERisk(data models.ScoreData) (riskValue uint64, err error)
	GetIdealAge(cveRiskValue uint64, data models.ScoreData) (ageMin, ageMax uint64, err error)
}
