package auth

import (
	"context"
	"fmt"

	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
	"github.com/cardio-analyst/backend/pkg/model"
)

func (c *Client) GetTokens(ctx context.Context, credentials model.Credentials, userIP string) (model.Tokens, error) {
	request := &pb.GetTokensRequest{
		Password: credentials.Password,
		Ip:       userIP,
	}

	if credentials.Login != "" {
		request.Login = &credentials.Login
	}
	if credentials.Email != "" {
		request.Email = &credentials.Email
	}

	response, err := c.client.GetTokens(ctx, request)
	if err != nil {
		return model.Tokens{}, err
	}

	if errorResponse := response.GetErrorResponse(); errorResponse != nil {
		switch errorResponse.GetErrorCode() {
		case pb.ErrorCode_INVALID_DATA:
			return model.Tokens{}, model.ErrInvalidCredentials
		case pb.ErrorCode_WRONG_CREDENTIALS:
			return model.Tokens{}, model.ErrWrongCredentials
		default:
			return model.Tokens{}, fmt.Errorf("unknown error code %v", errorResponse.GetErrorCode().String())
		}
	}

	successResponse := response.GetSuccessResponse()

	return model.Tokens{
		AccessToken:  successResponse.GetAccessToken(),
		RefreshToken: successResponse.GetRefreshToken(),
	}, nil
}

func (c *Client) RefreshTokens(ctx context.Context, refreshToken, userIP string) (model.Tokens, error) {
	request := &pb.RefreshTokensRequest{
		RefreshToken: refreshToken,
		Ip:           userIP,
	}

	response, err := c.client.RefreshTokens(ctx, request)
	if err != nil {
		return model.Tokens{}, err
	}

	if errorResponse := response.GetErrorResponse(); errorResponse != nil {
		switch errorResponse.GetErrorCode() {
		case pb.ErrorCode_WRONG_REFRESH_TOKEN:
			return model.Tokens{}, model.ErrWrongToken
		case pb.ErrorCode_REFRESH_TOKEN_EXPIRED:
			return model.Tokens{}, model.ErrTokenIsExpired
		case pb.ErrorCode_IP_NOT_ALLOWED:
			return model.Tokens{}, model.ErrIPIsNotInWhitelist
		default:
			return model.Tokens{}, fmt.Errorf("unknown error code %v", errorResponse.GetErrorCode().String())
		}
	}

	successResponse := response.GetSuccessResponse()

	return model.Tokens{
		AccessToken:  successResponse.GetAccessToken(),
		RefreshToken: successResponse.GetRefreshToken(),
	}, nil
}
