package users

import (
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

var _ service.UserService = (*userService)(nil)

type userService struct {
	db storage.UserStorage
}

func NewUserService(db storage.UserStorage) *userService {
	return &userService{db: db}
}

func (s *userService) Create(firstName, lastName, middleName, region, login, password string) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
