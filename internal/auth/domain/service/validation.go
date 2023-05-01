package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const (
	minPasswordLength = 7
	maxPasswordLength = 255
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (s *ValidationService) ValidateUser(user model.User) error {
	var checkPassword bool
	if user.Password != "" {
		checkPassword = true
	}

	err := validation.ValidateStruct(&user,
		validation.Field(&user.Role, validation.Required, validation.In(model.UserRoleCustomer, model.UserRoleModerator)),
		validation.Field(&user.FirstName, validation.Required),
		validation.Field(&user.LastName, validation.Required),
		validation.Field(&user.Region, validation.When(user.Role == model.UserRoleCustomer, validation.Required)),
		validation.Field(&user.BirthDate, validation.When(user.Role == model.UserRoleCustomer, validation.By(func(value any) error {
			date, ok := value.(model.Date)
			if !ok {
				return errors.New("cannot cast to date")
			}
			return s.ValidateDate(date)
		}))),
		validation.Field(&user.Login, validation.Required, validation.Match(regexp.MustCompile("^[^@]+$"))),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.When(
			checkPassword,
			validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		),
		validation.Field(&user.SecretKey, validation.When(user.Role == model.UserRoleModerator, validation.Required)),
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

		if validationError, found := validationErrors["role"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidFirstName, validationError)
		}
		if validationError, found := validationErrors["firstName"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidFirstName, validationError)
		}
		if validationError, found := validationErrors["lastName"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidLastName, validationError)
		}
		if validationError, found := validationErrors["region"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidRegion, validationError)
		}
		if validationError, found := validationErrors["birthDate"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidBirthDate, validationError)
		}
		if validationError, found := validationErrors["login"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidLogin, validationError)
		}
		if validationError, found := validationErrors["email"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidEmail, validationError)
		}
		if validationError, found := validationErrors["password"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidPassword, validationError)
		}
		if validationError, found := validationErrors["secretKey"]; found {
			return fmt.Errorf("%w: %v", model.ErrInvalidSecretKey, validationError)
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
		validation.Field(&credentials.Login, validation.When(credentials.Email == "", validation.Required)),
		validation.Field(&credentials.Email, validation.When(credentials.Login == "", validation.Required, is.Email)),
		validation.Field(&credentials.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
	)
}
