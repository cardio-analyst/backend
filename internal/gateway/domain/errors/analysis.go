package errors

import "errors"

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
