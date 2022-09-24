package service

type UserService interface {
	Create(firstName, lastName, middleName, region, login, password string) (userID uint64, err error)
}
