package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidRole      = errors.New("invalid role")
	ErrInvalidFirstName = errors.New("invalid first name")
	ErrInvalidLastName  = errors.New("invalid last name")
	ErrInvalidRegion    = errors.New("invalid region")
	ErrInvalidBirthDate = errors.New("invalid birth date")
	ErrInvalidLogin     = errors.New("invalid login")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidUserData  = errors.New("invalid user data")

	ErrUserLoginAlreadyOccupied = errors.New("user with such login is already registered")
	ErrUserEmailAlreadyOccupied = errors.New("user with such email is already registered")

	ErrUserNotFound = errors.New("user not found")
)

type UserRole string

const (
	UserRoleCustomer      UserRole = "CUSTOMER"
	UserRoleModerator     UserRole = "MODERATOR"
	UserRoleAdministrator UserRole = "ADMINISTRATOR"
)

type User struct {
	ID         uint64   `bson:"id" json:"-"`
	Role       UserRole `bson:"role" json:"-"`
	Login      string   `bson:"login" json:"login"`
	Email      string   `bson:"email" json:"email"`
	FirstName  string   `bson:"first_name" json:"firstName"`
	LastName   string   `bson:"last_name" json:"lastName"`
	Password   string   `bson:"password_hash,omitempty" json:"password,omitempty"`
	MiddleName string   `bson:"middle_name,omitempty" json:"middleName,omitempty"`
	Region     string   `bson:"region,omitempty" json:"region,omitempty"`
	BirthDate  Date     `bson:"birth_date,omitempty" json:"birthDate,omitempty"`
	SecretKey  string   `bson:"secret_key,omitempty" json:"secretKey,omitempty"`
}

func (u User) Age() int {
	today := time.Now().In(u.BirthDate.Location())

	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)

	by, bm, bd := u.BirthDate.Date()
	birthDate := time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)

	if today.Before(birthDate) {
		return 0
	}

	age := ty - by

	anniversary := birthDate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}

	return age
}

type CriteriaSeparator string

// criteria separators are used in query statement between arguments
const (
	CriteriaSeparatorAND CriteriaSeparator = "AND"
	CriteriaSeparatorOR  CriteriaSeparator = "OR"
)

type UserCriteria struct {
	ID                uint64
	Login             string
	Email             string
	PasswordHash      string
	Limit             int64
	Page              int64
	CriteriaSeparator CriteriaSeparator // required, takes value of CriteriaSeparatorAND or CriteriaSeparatorOR
}
