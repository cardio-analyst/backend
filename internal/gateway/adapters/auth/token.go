package auth

import (
	"context"
	"fmt"

	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
	"github.com/cardio-analyst/backend/pkg/model"
)

func (c *Client) GetTokens(ctx context.Context, credentials model.Credentials, userIP string, userRole model.UserRole) (model.Tokens, error) {
	role, err := userRolePB(userRole)
	if err != nil {
		return model.Tokens{}, err
	}

	request := &pb.GetTokensRequest{
		Password: credentials.Password,
		Ip:       userIP,
		UserRole: role,
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
		case pb.ErrorCode_FORBIDDEN_BY_ROLE:
			return model.Tokens{}, model.ErrForbiddenByRole
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

func (c *Client) RefreshTokens(ctx context.Context, refreshToken, userIP string, userRole model.UserRole) (model.Tokens, error) {
	role, err := userRolePB(userRole)
	if err != nil {
		return model.Tokens{}, err
	}

	request := &pb.RefreshTokensRequest{
		RefreshToken: refreshToken,
		Ip:           userIP,
		UserRole:     role,
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
		case pb.ErrorCode_FORBIDDEN_BY_ROLE:
			return model.Tokens{}, model.ErrForbiddenByRole
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

func userRolePB(role model.UserRole) (pb.UserRole, error) {
	switch role {
	case model.UserRoleCustomer:
		return pb.UserRole_CUSTOMER, nil
	case model.UserRoleModerator:
		return pb.UserRole_MODERATOR, nil
	case model.UserRoleAdministrator:
		return pb.UserRole_ADMINISTRATOR, nil
	default:
		return pb.UserRole(-1), fmt.Errorf("undefined user role: %q", role)
	}
}
