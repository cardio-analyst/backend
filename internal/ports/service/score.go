package service

import "github.com/cardio-analyst/backend/internal/domain/models"

type ScoreService interface {
	GetCveRisk(cveRiskData models.CveRiskData) (cveRisk uint64, err error)
}
