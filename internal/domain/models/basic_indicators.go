package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/cardio-analyst/backend/internal/domain/common"
)

type BasicIndicators struct {
	ID                     uint64   `json:"id,omitempty" db:"id"`
	UserID                 uint64   `json:"-" db:"user_id"`
	Weight                 *float64 `json:"weight" db:"weight"`
	Height                 *float64 `json:"height" db:"height"`
	BodyMassIndex          *float64 `json:"bodyMassIndex" db:"body_mass_index"`
	WaistSize              *float64 `json:"waistSize" db:"waist_size"`
	Gender                 *string  `json:"gender" db:"gender"`
	SBPLevel               *float64 `json:"sbpLevel" db:"sbp_level"`
	Smoking                *bool    `json:"smoking" db:"smoking"`
	TotalCholesterolLevel  *float64 `json:"totalCholesterolLevel" db:"total_cholesterol_level"`
	CVEventsRiskValue      *int64   `json:"cvEventsRiskValue" db:"cv_events_risk_value"`
	IdealCardiovascularAge *int64   `json:"idealCardiovascularAge" db:"ideal_cardiovascular_age"`
	CreatedAt              Datetime `json:"createdAt" db:"created_at"`
}

func (a BasicIndicators) Validate(updating bool) error {
	return validation.ValidateStruct(&a,
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
			validation.Required, validation.Min(16.0), validation.Max(60.0),
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
		validation.Field(&a.IdealCardiovascularAge, validation.When(
			a.IdealCardiovascularAge != nil,
			validation.Required, validation.Min(40), validation.Max(100),
		)),
	)
}
