package models

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// criteria separators are used in WHERE statement between arguments
const (
	CriteriaSeparatorAND = "AND"
	CriteriaSeparatorOR  = "OR"
)

// User TODO
type User struct {
	ID         uint64 `json:"-" db:"id"`
	FirstName  string `json:"firstName" db:"first_name"`
	LastName   string `json:"lastName" db:"last_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	Region     string `json:"region" db:"region"`
	BirthDate  Date   `json:"birthDate" db:"birth_date"`
	Login      string `json:"login" db:"login"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password_hash"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Region, validation.Required),
		validation.Field(&u.BirthDate, validation.By(u.BirthDate.Validate)),
		validation.Field(&u.Login, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(7, 255)),
	)
}

// UserCriteria TODO
type UserCriteria struct {
	Login             *string
	Email             *string
	PasswordHash      *string
	CriteriaSeparator string // required, "AND" or "OR"
}

// GetWhereStmtAndArgs TODO
func (c UserCriteria) GetWhereStmtAndArgs() (string, []interface{}) {
	whereStmtParts := make([]string, 0, 3)
	whereStmtArgs := make([]interface{}, 0, 3)
	currArgNum := 1

	if c.Login != nil {
		whereStmtParts = append(whereStmtParts, fmt.Sprintf("login=$%v", currArgNum))
		whereStmtArgs = append(whereStmtArgs, *c.Login)
		currArgNum++
	}
	if c.Email != nil {
		whereStmtParts = append(whereStmtParts, fmt.Sprintf("email=$%v", currArgNum))
		whereStmtArgs = append(whereStmtArgs, *c.Email)
		currArgNum++
	}
	if c.PasswordHash != nil {
		whereStmtParts = append(whereStmtParts, fmt.Sprintf("password_hash=$%v", currArgNum))
		whereStmtArgs = append(whereStmtArgs, *c.PasswordHash)
		currArgNum++
	}

	whereStmt := strings.Join(whereStmtParts, fmt.Sprintf(" %v ", c.CriteriaSeparator))

	return whereStmt, whereStmtArgs
}

// Credentials TODO
type Credentials struct {
	Login    *string
	Email    *string
	Password string
}
