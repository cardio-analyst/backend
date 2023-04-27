package grpc

import (
	"context"

	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
)

func (s *Server) GenerateSecretKey(_ context.Context, request *pb.GenerateSecretKeyRequest) (*pb.GenerateSecretKeyResponse, error) {
	secretKey, err := s.services.Auth().GenerateSecretKey(request.GetUserLogin(), request.GetUserEmail())
	if err != nil {
		return nil, err
	}
	return generateSecretKeySuccessResponse(secretKey), nil
}

func generateSecretKeySuccessResponse(secretKey string) *pb.GenerateSecretKeyResponse {
	return &pb.GenerateSecretKeyResponse{
		SecretKey: secretKey,
	}
}
