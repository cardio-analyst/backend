package grpc

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/cardio-analyst/backend/api/proto/analytics"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (s *Server) FindAllFeedbacks(_ context.Context, _ *emptypb.Empty) (*pb.FindAllFeedbacksResponse, error) {
	feedbacks, err := s.services.Feedback().FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Feedback, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		var userMiddleName *string
		if middleName := feedback.UserMiddleName; middleName != "" {
			userMiddleName = &middleName
		}

		var textMessage *string
		if message := feedback.Message; message != "" {
			textMessage = &message
		}

		createdAt := timestamppb.New(feedback.CreatedAt.Time)

		result = append(result, &pb.Feedback{
			Id:             feedback.ID,
			UserId:         feedback.UserID,
			UserFirstName:  feedback.UserFirstName,
			UserLastName:   feedback.UserLastName,
			UserMiddleName: userMiddleName,
			UserLogin:      feedback.UserLogin,
			UserEmail:      feedback.UserEmail,
			Mark:           int32(feedback.Mark),
			Message:        textMessage,
			Version:        feedback.Version,
			Viewed:         feedback.Viewed,
			CreatedAt:      createdAt,
		})
	}

	return &pb.FindAllFeedbacksResponse{
		Feedbacks: result,
	}, nil
}

func (s *Server) ToggleFeedbackViewed(_ context.Context, request *pb.ToggleFeedbackViewedRequest) (*pb.ToggleFeedbackViewedResponse, error) {
	if err := s.services.Feedback().ToggleFeedbackViewed(request.GetId()); err != nil {
		if errors.Is(err, model.ErrFeedbackNotFound) {
			return toggleFeedbackViewedErrorResponse(pb.ErrorCode_FEEDBACK_NOT_FOUND), nil
		}
		return nil, err
	}
	return toggleFeedbackViewedSuccessResponse(), nil
}

func toggleFeedbackViewedSuccessResponse() *pb.ToggleFeedbackViewedResponse {
	return &pb.ToggleFeedbackViewedResponse{
		Response: &pb.ToggleFeedbackViewedResponse_SuccessResponse{
			SuccessResponse: &emptypb.Empty{},
		},
	}
}

func toggleFeedbackViewedErrorResponse(errorCode pb.ErrorCode) *pb.ToggleFeedbackViewedResponse {
	return &pb.ToggleFeedbackViewedResponse{
		Response: &pb.ToggleFeedbackViewedResponse_ErrorResponse{
			ErrorResponse: &pb.ErrorResponse{
				ErrorCode: errorCode,
			},
		},
	}
}
