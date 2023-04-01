package errors

import "errors"

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
