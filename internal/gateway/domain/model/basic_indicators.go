package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cardio-analyst/backend/internal/gateway/domain/common"
	"github.com/cardio-analyst/backend/internal/pkg/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidWeight                       = errors.New("invalid weight value")
	ErrInvalidHeight                       = errors.New("invalid height value")
	ErrInvalidBodyMassIndex                = errors.New("invalid bodyMassIndex value")
	ErrInvalidWaistSize                    = errors.New("invalid waistSize value")
	ErrInvalidGender                       = errors.New("invalid gender value")
	ErrInvalidSBPLevel                     = errors.New("invalid sbpLevel value")
	ErrInvalidTotalCholesterolLevel        = errors.New("invalid totalCholesterolLevel value")
	ErrInvalidCVEventsRiskValue            = errors.New("invalid cvEventsRiskValue value")
	ErrInvalidIdealCardiovascularAgesRange = errors.New("invalid idealCardiovascularAgesRange value")
	ErrInvalidBasicIndicatorsData          = errors.New("invalid basic indicators data")
	ErrBasicIndicatorsRecordNotFound       = errors.New("basic indicators record with this id not found")
)

type BasicIndicators struct {
	ID                           uint64         `json:"id,omitempty" db:"id"`
	UserID                       uint64         `json:"-" db:"user_id"`
	Weight                       *float64       `json:"weight" db:"weight"`
	Height                       *float64       `json:"height" db:"height"`
	BodyMassIndex                *float64       `json:"bodyMassIndex" db:"body_mass_index"`
	WaistSize                    *float64       `json:"waistSize" db:"waist_size"`
	Gender                       *string        `json:"gender" db:"gender"`
	SBPLevel                     *float64       `json:"sbpLevel" db:"sbp_level"`
	Smoking                      *bool          `json:"smoking" db:"smoking"`
	TotalCholesterolLevel        *float64       `json:"totalCholesterolLevel" db:"total_cholesterol_level"`
	CVEventsRiskValue            *int64         `json:"cvEventsRiskValue" db:"cv_events_risk_value"`
	IdealCardiovascularAgesRange *string        `json:"idealCardiovascularAgesRange" db:"ideal_cardiovascular_ages_range"`
	Scale                        string         `json:"scale" db:"-"`
	CreatedAt                    model.Datetime `json:"createdAt" db:"created_at"`
}

func (a BasicIndicators) Validate(updating bool) error {
	err := validation.ValidateStruct(&a,
		validation.Field(&a.ID, validation.When(
			updating,
			validation.Required,
		)),
		validation.Field(&a.UserID, validation.Required),
		validation.Field(&a.Weight, validation.When(
			a.Weight != nil,
			validation.Required, validation.Min(40.0), validation.Max(160.0),
		)),
		validation.Field(&a.Height, validation.When(
			a.Height != nil,
			validation.Required, validation.Min(145.0), validation.Max(240.0),
		)),
		validation.Field(&a.BodyMassIndex, validation.When(
			a.BodyMassIndex != nil,
			validation.Required, validation.Min(1.0), validation.Max(60.0),
		)),
		validation.Field(&a.WaistSize, validation.When(
			a.WaistSize != nil,
			validation.Required, validation.Min(50.0), validation.Max(190.0),
		)),
		validation.Field(&a.Gender, validation.When(
			a.Gender != nil,
			validation.Required, validation.In(common.UserGenderMale, common.UserGenderFemale, common.UserGenderUnknown),
		)),
		validation.Field(&a.SBPLevel, validation.When(
			a.SBPLevel != nil,
			validation.Required, validation.Min(80.0), validation.Max(250.0),
		)),
		validation.Field(&a.TotalCholesterolLevel, validation.When(
			a.TotalCholesterolLevel != nil,
			validation.Required, validation.Min(3.0), validation.Max(15.2),
		)),
		validation.Field(&a.CVEventsRiskValue, validation.When(
			a.CVEventsRiskValue != nil,
			validation.Min(0), validation.Max(100),
		)),
		validation.Field(&a.IdealCardiovascularAgesRange, validation.When(
			a.IdealCardiovascularAgesRange != nil,
			validation.Required,
		)),
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

		if validationError, found := validationErrors["weight"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidWeight, validationError)
		}
		if validationError, found := validationErrors["height"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidHeight, validationError)
		}
		if validationError, found := validationErrors["bodyMassIndex"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidBodyMassIndex, validationError)
		}
		if validationError, found := validationErrors["waistSize"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidWaistSize, validationError)
		}
		if validationError, found := validationErrors["gender"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidGender, validationError)
		}
		if validationError, found := validationErrors["sbpLevel"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidSBPLevel, validationError)
		}
		if validationError, found := validationErrors["totalCholesterolLevel"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidTotalCholesterolLevel, validationError)
		}
		if validationError, found := validationErrors["cvEventsRiskValue"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidCVEventsRiskValue, validationError)
		}
		if validationError, found := validationErrors["idealCardiovascularAgesRange"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidIdealCardiovascularAgesRange, validationError)
		}

		return ErrInvalidBasicIndicatorsData
	}
	return nil
}
