package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/cardio-analyst/backend/api/proto/analytics"
)

func (s *Server) FindAllFeedbacks(_ context.Context, _ *emptypb.Empty) (*pb.FindAllFeedbacksResponse, error) {
	feedbacks, err := s.services.Feedback().FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Feedback, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		var userMiddleName *string
		if feedback.UserMiddleName != "" {
			userMiddleName = &feedback.UserMiddleName
		}

		var message *string
		if feedback.Message != "" {
			message = &feedback.Message
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
			Message:        message,
			CreatedAt:      createdAt,
		})
	}

	return &pb.FindAllFeedbacksResponse{
		Feedbacks: result,
	}, nil
}
