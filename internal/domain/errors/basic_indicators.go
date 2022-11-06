package errors

import "errors"

var (
	ErrInvalidBasicIndicatorsData    = errors.New("invalid basic indicators data")
	ErrBasicIndicatorsRecordNotFound = errors.New("basic indicators record with this id not found")
)
