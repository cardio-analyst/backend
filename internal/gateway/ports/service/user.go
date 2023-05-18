package service

import (
	"context"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type UserService interface {
	GetOne(ctx context.Context, criteria model.UserCriteria) (user model.User, err error)
	GetList(ctx context.Context, criteria model.UserCriteria) (users []model.User, totalPages int64, err error)
	Update(ctx context.Context, userData model.User) (err error)
}
