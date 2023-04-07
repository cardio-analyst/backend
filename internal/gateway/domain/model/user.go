package model

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/common"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserCredentials struct {
	LoginOrEmail string `json:"loginOrEmail"`
	Password     string `json:"password"`
}

func (r UserCredentials) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.LoginOrEmail, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(common.MinPasswordLength, common.MaxPasswordLength)),
	)
}
