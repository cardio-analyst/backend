package service

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type AnalysisService interface {
	Create(analysisData domain.Analysis) (err error)
	Update(analysisData domain.Analysis) (err error)
	FindAll(userID uint64) (analysisDataList []*domain.Analysis, err error)
}
