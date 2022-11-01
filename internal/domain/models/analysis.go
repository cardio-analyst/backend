package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Analysis struct {
	ID                              uint64   `json:"id,omitempty" db:"id"`
	UserID                          uint64   `json:"-" db:"user_id"`
	HighDensityCholesterol          *float64 `json:"highDensityCholesterol" db:"high_density_cholesterol"`
	LowDensityCholesterol           *float64 `json:"lowDensityCholesterol" db:"low_density_cholesterol"`
	Triglycerides                   *float64 `json:"triglycerides" db:"triglycerides"`
	Lipoprotein                     *float64 `json:"lipoprotein" db:"lipoprotein"`
	HighlySensitiveCReactiveProtein *float64 `json:"highlySensitiveCReactiveProtein" db:"highly_sensitive_c_reactive_protein"`
	AtherogenicityCoefficient       *float64 `json:"atherogenicityCoefficient" db:"atherogenicity_coefficient"`
	Creatinine                      *float64 `json:"creatinine" db:"creatinine"`
	AtheroscleroticPlaquesPresence  *bool    `json:"atheroscleroticPlaquesPresence" db:"atherosclerotic_plaques_presence"`
	CreatedAt                       Datetime `json:"createdAt" db:"created_at"`
}

func (a Analysis) Validate(updating bool) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ID, validation.When(
			updating,
			validation.Required,
		)),
		validation.Field(&a.UserID, validation.Required),
		validation.Field(&a.HighDensityCholesterol, validation.When(
			a.HighDensityCholesterol != nil,
			validation.Required, validation.Min(0.5), validation.Max(5.5),
		)),
		validation.Field(&a.LowDensityCholesterol, validation.When(
			a.LowDensityCholesterol != nil,
			validation.Required, validation.Min(0.5), validation.Max(8.5),
		)),
		validation.Field(&a.Triglycerides, validation.When(
			a.Triglycerides != nil,
			validation.Required, validation.Min(0.2), validation.Max(8.5),
		)),
		validation.Field(&a.Lipoprotein, validation.When(
			a.Lipoprotein != nil,
			validation.Min(0.0), validation.Max(10.0),
		)),
		validation.Field(&a.HighlySensitiveCReactiveProtein, validation.When(
			a.HighlySensitiveCReactiveProtein != nil,
			validation.Required, validation.Min(0.1), validation.Max(12.0),
		)),
		validation.Field(&a.AtherogenicityCoefficient, validation.When(
			a.AtherogenicityCoefficient != nil,
			validation.Required, validation.Min(0.1), validation.Max(8.0),
		)),
		validation.Field(&a.Creatinine, validation.When(
			a.Creatinine != nil,
			validation.Required, validation.Min(20.0), validation.Max(500.0),
		)),
	)
}
