package client

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type Analytics interface {
	FindAllFeedbacks(ctx context.Context, criteria model.FeedbackCriteria) (feedbacks []model.Feedback, totalPages int64, err error)
	ToggleFeedbackViewed(ctx context.Context, id uint64) error

	UsersByRegions(ctx context.Context) (map[string]int64, error)
}
