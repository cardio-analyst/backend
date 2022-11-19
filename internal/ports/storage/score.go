package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

type ScoreRepository interface {
	GetCveRisk(cveRiskData models.CveRiskData) (cveRisk uint64, err error)
}
