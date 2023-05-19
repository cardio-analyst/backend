package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/cardio-analyst/backend/api/proto/analytics"
)

func (s *Server) UsersByRegions(_ context.Context, _ *emptypb.Empty) (*pb.UsersByRegionsResponse, error) {
	usersByRegions, err := s.services.Statistics().AllUsersByRegions()
	if err != nil {
		return nil, err
	}

	return &pb.UsersByRegionsResponse{
		UsersByRegions: usersByRegions,
	}, nil
}
