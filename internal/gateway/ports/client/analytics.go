package client

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type Analytics interface {
	FindAllFeedbacks(ctx context.Context) ([]model.Feedback, error)
	ToggleFeedbackViewed(ctx context.Context, id uint64) error
}
