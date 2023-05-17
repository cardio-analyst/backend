package analytics

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/cardio-analyst/backend/api/proto/analytics"
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
			Version:        feedback.GetVersion(),
			Viewed:         feedback.GetViewed(),
			CreatedAt: model.Datetime{
				Time: feedback.GetCreatedAt().AsTime(),
			},
		})
	}

	return feedbacks, nil
}

func (c *Client) ToggleFeedbackViewed(ctx context.Context, id uint64) error {
	request := &pb.ToggleFeedbackViewedRequest{
		Id: id,
	}

	response, err := c.client.ToggleFeedbackViewed(ctx, request)
	if err != nil {
		return err
	}

	if errorResponse := response.GetErrorResponse(); errorResponse != nil {
		switch errorResponse.GetErrorCode() {
		case pb.ErrorCode_FEEDBACK_NOT_FOUND:
			return model.ErrFeedbackNotFound
		default:
			return fmt.Errorf("unknown error code %v", errorResponse.GetErrorCode().String())
		}
	}

	return nil
}
