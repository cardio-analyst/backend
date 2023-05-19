package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type StatisticsService interface {
	NotifyRegistration(user model.User)
	UsersByRegions(ctx context.Context) (map[string]int64, error)
}
