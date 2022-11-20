package service

import "github.com/cardio-analyst/backend/internal/domain/models"

type ScoreService interface {
	GetCVERisk(cveRiskData models.CVERiskData) (riskValue uint64, err error)
}
