package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

type ScoreRepository interface {
	GetCVERisk(cveRiskData models.CVERiskData) (riskValue uint64, err error)
}
