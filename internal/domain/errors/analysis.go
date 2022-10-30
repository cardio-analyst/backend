package errors

import "errors"

var (
	ErrInvalidAnalysisData    = errors.New("invalid analysis data")
	ErrAnalysisRecordNotFound = errors.New("analysis record with this id not found")
)
