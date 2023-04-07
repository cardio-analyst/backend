package grpc

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	domain "github.com/cardio-analyst/backend/internal/auth/domain/model"
	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
	"github.com/cardio-analyst/backend/pkg/model"
)

func (s *Server) GetTokens(ctx context.Context, request *pb.GetTokensRequest) (*pb.TokensResponse, error) {
	credentials := model.Credentials{
		Login:    request.GetLogin(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	if err := s.services.Validation().ValidateCredentials(credentials); err != nil {
		log.Errorf("validating user credentials: %v", err)
		return tokensErrorResponse(pb.ErrorCode_INVALID_DATA), nil
	}

	tokens, err := s.services.Auth().GetTokens(ctx, credentials, request.GetIp())
	if err != nil {
		log.Errorf("receiving tokens from auth service: %v", err)
		if errors.Is(err, model.ErrWrongCredentials) {
			return tokensErrorResponse(pb.ErrorCode_WRONG_CREDENTIALS), nil
		}
		return nil, err
	}

	log.Debugf("get tokens payload: access: %q, refresh: %q", tokens.AccessToken, tokens.RefreshToken)

	return tokensSuccessResponse(tokens.AccessToken, tokens.RefreshToken), nil
}

func (s *Server) RefreshTokens(ctx context.Context, request *pb.RefreshTokensRequest) (*pb.TokensResponse, error) {
	tokens, err := s.services.Auth().RefreshTokens(ctx, request.GetRefreshToken(), request.GetIp())
	if err != nil {
		log.Errorf("refreshing tokens in auth service: %v", err)
		switch {
		case errors.Is(err, model.ErrWrongToken) || errors.Is(err, domain.ErrSessionNotFound):
			return tokensErrorResponse(pb.ErrorCode_WRONG_REFRESH_TOKEN), nil
		case errors.Is(err, model.ErrTokenIsExpired):
			return tokensErrorResponse(pb.ErrorCode_REFRESH_TOKEN_EXPIRED), nil
		case errors.Is(err, model.ErrIPIsNotInWhitelist):
			return tokensErrorResponse(pb.ErrorCode_IP_NOT_ALLOWED), nil
		}
		return nil, err
	}

	log.Debugf("refresh tokens payload: access: %q, refresh: %q", tokens.AccessToken, tokens.RefreshToken)

	return tokensSuccessResponse(tokens.AccessToken, tokens.RefreshToken), nil
}

func tokensSuccessResponse(accessToken, refreshToken string) *pb.TokensResponse {
	return &pb.TokensResponse{
		Response: &pb.TokensResponse_SuccessResponse{
			SuccessResponse: &pb.TokensSuccessResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	}
}

func tokensErrorResponse(errorCode pb.ErrorCode) *pb.TokensResponse {
	return &pb.TokensResponse{
		Response: &pb.TokensResponse_ErrorResponse{
			ErrorResponse: &pb.ErrorResponse{
				ErrorCode: errorCode,
			},
		},
	}
}
