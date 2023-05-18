package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/cardio-analyst/backend/api/proto/analytics"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (s *Server) FindAllFeedbacks(_ context.Context, request *pb.FindAllFeedbacksRequest) (*pb.FindAllFeedbacksResponse, error) {
	criteria := model.FeedbackCriteria{
		Viewed: request.Viewed,
		Limit:  request.GetLimit(),
		Page:   request.GetPage(),
	}

	switch request.GetMarkOrdering() {
	case pb.OrderingType_DISABLED:
		criteria.MarkOrdering = model.OrderingTypeDisabled
	case pb.OrderingType_ASCENDING:
		criteria.MarkOrdering = model.OrderingTypeASC
	case pb.OrderingType_DESCENDING:
		criteria.MarkOrdering = model.OrderingTypeDESC
	default:
		return nil, fmt.Errorf("unknown mark ordering type: %v", request.GetMarkOrdering())
	}

	switch request.GetVersionOrdering() {
	case pb.OrderingType_DISABLED:
		criteria.VersionOrdering = model.OrderingTypeDisabled
	case pb.OrderingType_ASCENDING:
		criteria.VersionOrdering = model.OrderingTypeASC
	case pb.OrderingType_DESCENDING:
		criteria.VersionOrdering = model.OrderingTypeDESC
	default:
		return nil, fmt.Errorf("unknown version ordering type: %v", request.GetVersionOrdering())
	}

	feedbacks, totalPages, err := s.services.Feedback().FindAll(criteria)
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
		Feedbacks:  result,
		TotalPages: totalPages,
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
