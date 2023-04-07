package model

import (
	"encoding/json"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidHighDensityCholesterol          = errors.New("invalid highDensityCholesterol value")
	ErrInvalidLowDensityCholesterol           = errors.New("invalid lowDensityCholesterol value")
	ErrInvalidTriglycerides                   = errors.New("invalid triglycerides value")
	ErrInvalidLipoprotein                     = errors.New("invalid lipoprotein value")
	ErrInvalidHighlySensitiveCReactiveProtein = errors.New("invalid highlySensitiveCReactiveProtein value")
	ErrInvalidAtherogenicityCoefficient       = errors.New("invalid atherogenicityCoefficient value")
	ErrInvalidCreatinine                      = errors.New("invalid creatinine value")
	ErrInvalidAnalysisData                    = errors.New("invalid data")
	ErrAnalysisRecordNotFound                 = errors.New("analysis record with this id not found")
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
	err := validation.ValidateStruct(&a,
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

		if validationError, found := validationErrors["highDensityCholesterol"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidHighDensityCholesterol, validationError)
		}
		if validationError, found := validationErrors["lowDensityCholesterol"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidLowDensityCholesterol, validationError)
		}
		if validationError, found := validationErrors["triglycerides"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidTriglycerides, validationError)
		}
		if validationError, found := validationErrors["lipoprotein"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidLipoprotein, validationError)
		}
		if validationError, found := validationErrors["highlySensitiveCReactiveProtein"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidHighlySensitiveCReactiveProtein, validationError)
		}
		if validationError, found := validationErrors["atherogenicityCoefficient"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidAtherogenicityCoefficient, validationError)
		}
		if validationError, found := validationErrors["creatinine"]; found {
			return fmt.Errorf("%w: %v", ErrInvalidCreatinine, validationError)
		}

		return ErrInvalidAnalysisData
	}
	return nil
}
