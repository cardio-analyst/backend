package service

import (
	"context"

	"github.com/cardio-analyst/backend/pkg/model"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (err error)
	GetOne(ctx context.Context, criteria model.UserCriteria) (user model.User, err error)
}
