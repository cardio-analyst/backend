package grpc

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
	"github.com/cardio-analyst/backend/pkg/model"
)

func (s *Server) SaveUser(ctx context.Context, request *pb.SaveUserRequest) (*pb.SaveUserResponse, error) {
	user := model.User{
		Role:       model.UserRole(request.GetRole().String()),
		Login:      request.GetLogin(),
		Email:      request.GetEmail(),
		FirstName:  request.GetFirstName(),
		LastName:   request.GetLastName(),
		Password:   request.GetPassword(),
		MiddleName: request.GetMiddleName(),
		Region:     request.GetRegion(),
		BirthDate:  model.Date{Time: request.BirthDate.AsTime()},
		SecretKey:  request.GetSecretKey(),
	}

	if err := s.services.Validation().ValidateUser(user); err != nil {
		log.Errorf("validating user data: %v", err)
		switch {
		case errors.Is(err, model.ErrInvalidRole):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_ROLE), nil
		case errors.Is(err, model.ErrInvalidFirstName):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_FIRST_NAME), nil
		case errors.Is(err, model.ErrInvalidLastName):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_LAST_NAME), nil
		case errors.Is(err, model.ErrInvalidRegion):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_REGION), nil
		case errors.Is(err, model.ErrInvalidBirthDate):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_BIRTH_DATE), nil
		case errors.Is(err, model.ErrInvalidLogin):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_LOGIN), nil
		case errors.Is(err, model.ErrInvalidEmail):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_EMAIL), nil
		case errors.Is(err, model.ErrInvalidPassword):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_PASSWORD), nil
		case errors.Is(err, model.ErrInvalidSecretKey):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_SECRET_KEY), nil
		case errors.Is(err, model.ErrInvalidUserData):
			return saveUserErrorResponse(pb.ErrorCode_INVALID_DATA), nil
		}
		return nil, err
	}

	if err := s.services.Auth().VerifySecretKey(user); err != nil {
		log.Errorf("verifying user secret key: %v", err)
		if errors.Is(err, model.ErrWrongSecretKey) {
			return saveUserErrorResponse(pb.ErrorCode_WRONG_SECRET_KEY), nil
		}
		return nil, err
	}

	if err := s.services.User().Save(ctx, user); err != nil {
		log.Errorf("saving user: %v", err)
		switch {
		case errors.Is(err, model.ErrUserLoginAlreadyOccupied):
			return saveUserErrorResponse(pb.ErrorCode_LOGIN_ALREADY_OCCUPIED), nil
		case errors.Is(err, model.ErrUserEmailAlreadyOccupied):
			return saveUserErrorResponse(pb.ErrorCode_EMAIL_ALREADY_OCCUPIED), nil
		}
		return nil, err
	}

	log.Debugf("user with login %q, email %q and role %q successfully created", user.Login, user.Email, user.Role)

	return saveUserSuccessResponse(), nil
}

func saveUserSuccessResponse() *pb.SaveUserResponse {
	return &pb.SaveUserResponse{
		Response: &pb.SaveUserResponse_SuccessResponse{
			SuccessResponse: &emptypb.Empty{},
		},
	}
}

func saveUserErrorResponse(errorCode pb.ErrorCode) *pb.SaveUserResponse {
	return &pb.SaveUserResponse{
		Response: &pb.SaveUserResponse_ErrorResponse{
			ErrorResponse: &pb.ErrorResponse{
				ErrorCode: errorCode,
			},
		},
	}
}

func (s *Server) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var criteria model.UserCriteria
	if request.Id != nil {
		criteria.ID = request.GetId()
	}
	if request.Login != nil {
		criteria.Login = request.GetLogin()
	}
	if request.Email != nil {
		criteria.Email = request.GetEmail()
	}

	user, err := s.services.User().GetOne(ctx, criteria)
	if err != nil {
		log.Errorf("receiving user with login %v and email %v: %v", user.Login, user.Email, err)
		if errors.Is(err, model.ErrUserNotFound) {
			return getUserErrorResponse(pb.ErrorCode_USER_NOT_FOUND), nil
		}
		return nil, err
	}

	return getUserSuccessResponse(user), nil
}

func getUserSuccessResponse(user model.User) *pb.GetUserResponse {
	var role pb.UserRole
	switch user.Role {
	case model.UserRoleCustomer:
		role = pb.UserRole_CUSTOMER
	case model.UserRoleModerator:
		role = pb.UserRole_MODERATOR
	case model.UserRoleAdministrator:
		role = pb.UserRole_ADMINISTRATOR
	}

	birthDate := timestamppb.New(user.BirthDate.Time)

	return &pb.GetUserResponse{
		Response: &pb.GetUserResponse_SuccessResponse{
			SuccessResponse: &pb.GetUserSuccessResponse{
				User: &pb.User{
					Id:           user.ID,
					Role:         role,
					Login:        user.Login,
					Email:        user.Email,
					FirstName:    user.FirstName,
					LastName:     user.LastName,
					PasswordHash: user.Password,
					MiddleName:   &user.MiddleName,
					Region:       &user.Region,
					BirthDate:    birthDate,
				},
			},
		},
	}
}

func getUserErrorResponse(errorCode pb.ErrorCode) *pb.GetUserResponse {
	return &pb.GetUserResponse{
		Response: &pb.GetUserResponse_ErrorResponse{
			ErrorResponse: &pb.ErrorResponse{
				ErrorCode: errorCode,
			},
		},
	}
}

func (s *Server) IdentifyUser(ctx context.Context, request *pb.IdentifyUserRequest) (*pb.IdentifyUserResponse, error) {
	userID, userRole, err := s.services.Auth().IdentifyUser(ctx, request.GetAccessToken())
	if err != nil {
		log.Errorf("identifying user: %v", err)
		switch {
		case errors.Is(err, model.ErrWrongToken):
			return identifyUserErrorResponse(pb.ErrorCode_WRONG_ACCESS_TOKEN), nil
		case errors.Is(err, model.ErrTokenIsExpired):
			return identifyUserErrorResponse(pb.ErrorCode_ACCESS_TOKEN_EXPIRED), nil
		}
		return nil, err
	}

	log.Debugf("user successfully identified with id %v and role %q", userID, userRole)

	return identifyUserSuccessResponse(userID, userRole), nil
}

func identifyUserSuccessResponse(userID uint64, userRole model.UserRole) *pb.IdentifyUserResponse {
	var role pb.UserRole
	switch userRole {
	case model.UserRoleCustomer:
		role = pb.UserRole_CUSTOMER
	case model.UserRoleModerator:
		role = pb.UserRole_MODERATOR
	case model.UserRoleAdministrator:
		role = pb.UserRole_ADMINISTRATOR
	}

	return &pb.IdentifyUserResponse{
		Response: &pb.IdentifyUserResponse_SuccessResponse{
			SuccessResponse: &pb.IdentifyUserSuccessResponse{
				UserId: userID,
				Role:   role,
			},
		},
	}
}

func identifyUserErrorResponse(errorCode pb.ErrorCode) *pb.IdentifyUserResponse {
	return &pb.IdentifyUserResponse{
		Response: &pb.IdentifyUserResponse_ErrorResponse{
			ErrorResponse: &pb.ErrorResponse{
				ErrorCode: errorCode,
			},
		},
	}
}
