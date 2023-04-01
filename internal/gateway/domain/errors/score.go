package errors

import "errors"

var (
	ErrInvalidAge       = errors.New("invalid age value")
	ErrInvalidScoreData = errors.New("invalid SCORE data")
)
