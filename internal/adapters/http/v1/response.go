package v1

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// possible response designations
const (
	resultRegistered = "Registered"
	resultCreated    = "Created"
	resultUpdated    = "Updated"
)

// possible errors designations
const (
	errorParseRequestData              = "ParseRequestError"
	errorInvalidRequestData            = "InvalidRequestData"
	errorLoginAlreadyOccupied          = "LoginAlreadyOccupied"
	errorEmailAlreadyOccupied          = "EmailAlreadyOccupied"
	errorInternal                      = "InternalError"
	errorWrongCredentials              = "WrongCredentials"
	errorWrongAccessToken              = "WrongAccessToken"
	errorWrongRefreshToken             = "WrongRefreshToken"
	errorAccessTokenExpired            = "AccessTokenExpired"
	errorRefreshTokenExpired           = "RefreshTokenExpired"
	errorIPNotAllowed                  = "IPNotAllowed"
	errorWrongAuthHeader               = "WrongAuthHeader"
	errorAnalysisRecordNotFound        = "AnalysisRecordNotFound"
	errorBasicIndicatorsRecordNotFound = "BasicIndicatorsRecordNotFound"
	errorNotEnoughInformation          = "NotEnoughInformation"
)

var errorDescriptions = map[string]string{
	errorParseRequestData:              "Запрос составлен некорректно",
	errorInvalidRequestData:            "Ошибка валидации данных",
	errorLoginAlreadyOccupied:          "Выбранный логин уже занят",
	errorEmailAlreadyOccupied:          "Выбранный E-mail уже занят",
	errorInternal:                      "Внутренняя ошибка сервера",
	errorWrongCredentials:              "Некорректные данные для входа",
	errorWrongAccessToken:              "Некорректный access-токен",
	errorWrongRefreshToken:             "Некорректный refresh-токен",
	errorAccessTokenExpired:            "Время жизни access-токена истекло",
	errorRefreshTokenExpired:           "Время жизни refresh-токена истекло",
	errorIPNotAllowed:                  "Неизвестное устройство",
	errorWrongAuthHeader:               "Некорректный заголовок авторизации",
	errorAnalysisRecordNotFound:        "Запись о лабораторных исследованиях не найдена",
	errorBasicIndicatorsRecordNotFound: "Запись о базовых показателях не найдена",
	errorNotEnoughInformation:          "Недостаточно данных для выполнения запроса",
}

type response struct {
	Result      string `json:"result,omitempty"`
	Error       string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
}

func newResult(result string) *response {
	return &response{
		Result: result,
	}
}

func newError(c echo.Context, err error, error string) *response {
	log.WithFields(log.Fields{
		"error":      err.Error(),
		"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
	}).Error("error occurred")

	return &response{
		Error:       error,
		Description: errorDescriptions[error],
	}
}
