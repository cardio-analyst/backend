package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/cardio-analyst/backend/pkg/model"
)

const (
	minPasswordLength = 7
	maxPasswordLength = 255
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (s *ValidationService) ValidateUser(user model.User, checkPassword bool) error {
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Role, validation.Required, validation.In(model.UserRoleCustomer, model.UserRoleModerator)),
		validation.Field(&user.FirstName, validation.Required),
		validation.Field(&user.LastName, validation.Required),
		validation.Field(&user.Region, validation.Required),
		validation.Field(&user.BirthDate, validation.By(func(value any) error {
			date, ok := value.(model.Date)
			if !ok {
				return errors.New("cannot cast to date")
			}
			return s.ValidateDate(date)
		})),
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.When(
			checkPassword,
			validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		),
	)
	if err != nil {
		var errBytes []byte
		errBytes, err = json.Marshal(err)
		if err != nil {
			return err
		}

		var validationErrors map[string]string
		if err = json.Unmarshal(errBytes, &validationErrors); err != nil {
			return err
		}

		if validationError, found := validationErrors["Role"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidFirstName, validationError)
		}
		if validationError, found := validationErrors["FirstName"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidFirstName, validationError)
		}
		if validationError, found := validationErrors["LastName"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidLastName, validationError)
		}
		if validationError, found := validationErrors["Region"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidRegion, validationError)
		}
		if validationError, found := validationErrors["BirthDate"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidBirthDate, validationError)
		}
		if validationError, found := validationErrors["Login"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidLogin, validationError)
		}
		if validationError, found := validationErrors["Email"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidEmail, validationError)
		}
		if validationError, found := validationErrors["Password"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidPassword, validationError)
		}

		return model.ErrInvalidUserData
	}
	return nil
}

func (s *ValidationService) ValidateDate(date model.Date) error {
	if err := validation.Validate(date.String(), validation.Required, validation.Date(model.DateLayout)); err != nil {
		return err
	}
	return validation.Validate(date.Time, validation.Required, validation.Max(time.Now()))
}

func (s *ValidationService) ValidateCredentials(credentials model.Credentials) error {
	return validation.ValidateStruct(&credentials,
		validation.Field(&credentials.Login, validation.Required),
		validation.Field(&credentials.Email, validation.Required, is.Email),
		validation.Field(&credentials.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
	)
}
