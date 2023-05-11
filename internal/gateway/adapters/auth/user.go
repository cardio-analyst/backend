package auth

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/cardio-analyst/backend/api/proto/auth"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (c *Client) SaveUser(ctx context.Context, user model.User) error {
	role, err := userRolePB(user.Role)
	if err != nil {
		return err
	}
	if role == pb.UserRole_ADMINISTRATOR {
		return errors.New("forbidden to create the administrator")
	}

	birthDate := timestamppb.New(user.BirthDate.Time)

	var middleName *string
	if user.MiddleName != "" {
		middleName = &user.MiddleName
	}

	var region *string
	if user.Region != "" {
		region = &user.Region
	}

	var secretKey *string
	if user.SecretKey != "" {
		secretKey = &user.SecretKey
	}

	request := &pb.SaveUserRequest{
		Id:         user.ID,
		Role:       role,
		Login:      user.Login,
		Email:      user.Email,
		Password:   user.Password,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: middleName,
		Region:     region,
		BirthDate:  birthDate,
		SecretKey:  secretKey,
	}

	response, err := c.client.SaveUser(ctx, request)
	if err != nil {
		return err
	}

	if errorResponse := response.GetErrorResponse(); errorResponse != nil {
		switch errorResponse.GetErrorCode() {
		case pb.ErrorCode_INVALID_ROLE:
			return model.ErrInvalidRole
		case pb.ErrorCode_INVALID_FIRST_NAME:
			return model.ErrInvalidFirstName
		case pb.ErrorCode_INVALID_LAST_NAME:
			return model.ErrInvalidLastName
		case pb.ErrorCode_INVALID_REGION:
			return model.ErrInvalidRegion
		case pb.ErrorCode_INVALID_BIRTH_DATE:
			return model.ErrInvalidBirthDate
		case pb.ErrorCode_INVALID_LOGIN:
			return model.ErrInvalidLogin
		case pb.ErrorCode_INVALID_EMAIL:
			return model.ErrInvalidEmail
		case pb.ErrorCode_INVALID_PASSWORD:
			return model.ErrInvalidPassword
		case pb.ErrorCode_INVALID_DATA:
			return model.ErrInvalidUserData
		case pb.ErrorCode_LOGIN_ALREADY_OCCUPIED:
			return model.ErrUserLoginAlreadyOccupied
		case pb.ErrorCode_EMAIL_ALREADY_OCCUPIED:
			return model.ErrUserEmailAlreadyOccupied
		case pb.ErrorCode_INVALID_SECRET_KEY:
			return model.ErrInvalidSecretKey
		case pb.ErrorCode_WRONG_SECRET_KEY:
			return model.ErrWrongSecretKey
		default:
			return fmt.Errorf("unknown error code %v", errorResponse.GetErrorCode().String())
		}
	}

	return nil
}

func (c *Client) GetUser(ctx context.Context, criteria model.UserCriteria) (model.User, error) {
	request := new(pb.GetUserRequest)
	if criteria.ID != 0 {
		request.Id = &criteria.ID
	}
	if criteria.Login != "" {
		request.Login = &criteria.Login
	}
	if criteria.Email != "" {
		request.Email = &criteria.Email
	}

	response, err := c.client.GetUser(ctx, request)
	if err != nil {
		return model.User{}, err
	}

	if errorResponse := response.GetErrorResponse(); errorResponse != nil {
		switch errorResponse.GetErrorCode() {
		case pb.ErrorCode_USER_NOT_FOUND:
			return model.User{}, model.ErrUserNotFound
		default:
			return model.User{}, fmt.Errorf("unknown error code %v", errorResponse.GetErrorCode().String())
		}
	}

	user := response.GetSuccessResponse().GetUser()

	role, err := pbUserRole(user.GetRole())
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:         user.GetId(),
		Role:       role,
		Login:      user.GetLogin(),
		Email:      user.GetEmail(),
		FirstName:  user.GetFirstName(),
		LastName:   user.GetLastName(),
		Password:   user.PasswordHash,
		MiddleName: user.GetMiddleName(),
		Region:     user.GetRegion(),
		BirthDate:  model.Date{Time: user.GetBirthDate().AsTime()},
	}, nil
}

func (c *Client) IdentifyUser(ctx context.Context, accessToken string) (uint64, model.UserRole, error) {
	request := &pb.IdentifyUserRequest{
		AccessToken: accessToken,
	}

	response, err := c.client.IdentifyUser(ctx, request)
	if err != nil {
		return 0, "", err
	}

	if errorResponse := response.GetErrorResponse(); errorResponse != nil {
		switch errorResponse.GetErrorCode() {
		case pb.ErrorCode_WRONG_ACCESS_TOKEN:
			return 0, "", model.ErrWrongToken
		case pb.ErrorCode_ACCESS_TOKEN_EXPIRED:
			return 0, "", model.ErrTokenIsExpired
		default:
			return 0, "", fmt.Errorf("unknown error code %v", errorResponse.GetErrorCode().String())
		}
	}

	successResponse := response.GetSuccessResponse()

	role, err := pbUserRole(successResponse.GetRole())
	if err != nil {
		return 0, "", err
	}

	return successResponse.UserId, role, nil
}

func pbUserRole(userRolePB pb.UserRole) (model.UserRole, error) {
	switch userRolePB {
	case pb.UserRole_CUSTOMER:
		return model.UserRoleCustomer, nil
	case pb.UserRole_MODERATOR:
		return model.UserRoleModerator, nil
	case pb.UserRole_ADMINISTRATOR:
		return model.UserRoleAdministrator, nil
	default:
		return "", fmt.Errorf("unknown user role: %q", userRolePB.String())
	}
}
