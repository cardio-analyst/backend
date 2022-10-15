package http

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// possible response designations
const (
	resultRegistered = "Registered"
	resultUpdated    = "Updated"
)

// possible errors designations
const (
	errorParseRequestData     = "ParseRequestError"
	errorInvalidRequestData   = "InvalidRequestData"
	errorLoginAlreadyOccupied = "LoginAlreadyOccupied"
	errorEmailAlreadyOccupied = "EmailAlreadyOccupied"
	errorInternal             = "InternalError"
	errorWrongCredentials     = "WrongCredentials"
	errorWrongAccessToken     = "WrongAccessToken"
	errorWrongRefreshToken    = "WrongRefreshToken"
	errorAccessTokenExpired   = "AccessTokenExpired"
	errorRefreshTokenExpired  = "RefreshTokenExpired"
	errorIPNotAllowed         = "IPNotAllowed"
	errorWrongAuthHeader      = "WrongAuthHeader"
)

var errorDescriptions = map[string]string{
	errorParseRequestData:     "Запрос составлен некорректно",
	errorInvalidRequestData:   "Ошибка валидации данных",
	errorLoginAlreadyOccupied: "Выбранный логин уже занят",
	errorEmailAlreadyOccupied: "Выбранный E-mail уже занят",
	errorInternal:             "Внутренняя ошибка сервера",
	errorWrongCredentials:     "Некорректные данные для входа",
	errorWrongAccessToken:     "Некорректный access-токен",
	errorWrongRefreshToken:    "Некорректный refresh-токен",
	errorAccessTokenExpired:   "Время жизни access-токена истекло",
	errorRefreshTokenExpired:  "Время жизни refresh-токена истекло",
	errorIPNotAllowed:         "Неизвестное устройство",
	errorWrongAuthHeader:      "Некорректный заголовок авторизации",
}

type Response struct {
	Result      string `json:"result,omitempty"`
	Error       string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
}

func NewResult(result string) *Response {
	return &Response{
		Result: result,
	}
}

func NewError(c echo.Context, err error, error string) *Response {
	log.WithFields(log.Fields{
		"error":      err.Error(),
		"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
	}).Error("error occurred")

	return &Response{
		Error:       error,
		Description: errorDescriptions[error],
	}
}
