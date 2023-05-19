package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type StatisticsService interface {
	NotifyRegistration(user model.User)
	UsersByRegions(ctx context.Context) (map[string]int64, error)
	DiseasesByUsers(region string, regionUsers map[uint64]bool) (map[string]int64, error)
	SBPByUsers(region string, regionUsers map[uint64]bool) (map[string]int64, error)
	IdealCardiovascularAgesRangesByUsers(region string, regionUsers map[uint64]bool) (map[string]int64, error)
}
