package analytics

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (c *Client) FindAllFeedbacks(ctx context.Context) ([]model.Feedback, error) {
	request := &emptypb.Empty{}

	response, err := c.client.FindAllFeedbacks(ctx, request)
	if err != nil {
		return nil, err
	}

	feedbacks := make([]model.Feedback, 0, len(response.GetFeedbacks()))
	for _, feedback := range response.GetFeedbacks() {
		feedbacks = append(feedbacks, model.Feedback{
			ID:             feedback.GetId(),
			UserID:         feedback.GetUserId(),
			UserFirstName:  feedback.GetUserFirstName(),
			UserLastName:   feedback.GetUserLastName(),
			UserMiddleName: feedback.GetUserMiddleName(),
			UserLogin:      feedback.GetUserLogin(),
			UserEmail:      feedback.GetUserEmail(),
			Mark:           int16(feedback.GetMark()),
			Message:        feedback.GetMessage(),
			CreatedAt: model.Datetime{
				Time: feedback.GetCreatedAt().AsTime(),
			},
		})
	}

	return feedbacks, nil
}
