package http

import (
	"fmt"
	"net/http"
)

// possible response designations
const (
	resultOK = "Ok"
)

type Response struct {
	Result string `json:"result"`
}

func NewOKResponse() (int, *Response) {
	return http.StatusOK, &Response{
		Result: resultOK,
	}
}

// possible errors designations
const (
	errorParseRequestData           = "ParseRequestDataError"
	errorInvalidRequestData         = "InvalidRequestData"
	errorAlreadyRegisteredWithLogin = "AlreadyRegisteredWithLogin"
	errorAlreadyRegisteredWithEmail = "AlreadyRegisteredWithEmail"
	errorInternal                   = "InternalError"
	errorForbidden                  = "Forbidden"
)

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

func NewParseRequestDataErrorResponse(err error) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorParseRequestData,
		Description: fmt.Sprintf("failed to parse request data: %v", err),
	}
}

func NewInvalidRequestDataResponse(err error) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorInvalidRequestData,
		Description: fmt.Sprintf("validation failed: %v", err),
	}
}

func NewAlreadyRegisteredWithLoginResponse(login string) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorAlreadyRegisteredWithLogin,
		Description: fmt.Sprintf("user with login '%v' already registered", login),
	}
}

func NewAlreadyRegisteredWithEmailResponse(email string) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Error:       errorAlreadyRegisteredWithEmail,
		Description: fmt.Sprintf("user with email '%v' already registered", email),
	}
}

func NewInternalErrorResponse(err error) (int, *ErrorResponse) {
	return http.StatusInternalServerError, &ErrorResponse{
		Error:       errorInternal,
		Description: err.Error(),
	}
}

func NewForbiddenResponse(err error) (int, *ErrorResponse) {
	return http.StatusForbidden, &ErrorResponse{
		Error:       errorForbidden,
		Description: err.Error(),
	}
}
