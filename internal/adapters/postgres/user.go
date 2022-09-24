package postgres

import "github.com/cardio-analyst/backend/internal/ports/storage"

const userTable = "users"

var _ storage.UserStorage = (*Database)(nil)

func (d *Database) Create(firstName, lastName, middleName, region, login, password string) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
