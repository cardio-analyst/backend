package storage

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type ScoreRepository interface {
	GetCVERisk(data domain.ScoreData) (riskValue uint64, err error)
	GetIdealAge(cveRiskValue uint64, data domain.ScoreData) (ageMin, ageMax uint64, err error)
}
