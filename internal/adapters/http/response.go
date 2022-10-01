package http

import (
	"fmt"
	"net/http"
)

// possible errors designations
const (
	errorParseRequestData           = "ParseRequestDataError"
	errorInvalidRequestData         = "InvalidRequestData"
	errorAlreadyRegisteredWithLogin = "AlreadyRegisteredWithLogin"
	errorAlreadyRegisteredWithEmail = "AlreadyRegisteredWithEmail"
	errorInternal                   = "InternalError"
)

// ErrorResponse TODO
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

// NewParseRequestDataErrorResponse TODO
func NewParseRequestDataErrorResponse(err error) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorParseRequestData,
		Description: fmt.Sprintf("failed to parse request data: %v", err),
	}
}

// NewInvalidRequestDataResponse TODO
func NewInvalidRequestDataResponse(err error) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorInvalidRequestData,
		Description: fmt.Sprintf("validation failed: %v", err),
	}
}

// NewAlreadyRegisteredWithLoginResponse TODO
func NewAlreadyRegisteredWithLoginResponse(login string) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorAlreadyRegisteredWithLogin,
		Description: fmt.Sprintf("user with login '%v' already registered", login),
	}
}

// NewAlreadyRegisteredWithEmailResponse TODO
func NewAlreadyRegisteredWithEmailResponse(email string) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorAlreadyRegisteredWithEmail,
		Description: fmt.Sprintf("user with email '%v' already registered", email),
	}
}

// NewInternalErrorResponse TODO
func NewInternalErrorResponse(err error) (int, *ErrorResponse) {
	return http.StatusInternalServerError, &ErrorResponse{
		Error:       errorInternal,
		Description: err.Error(),
	}
}
