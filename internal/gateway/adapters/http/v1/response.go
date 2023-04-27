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
	resultEmailSent  = "Sent"
)

// possible errors designations
const (
	errorParseRequestData     = "ParseRequestError"
	errorInvalidRequestData   = "InvalidRequestData"
	errorLoginAlreadyOccupied = "LoginAlreadyOccupied"
	errorEmailAlreadyOccupied = "EmailAlreadyOccupied"
	errorInternal             = "InternalError"
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
	errorInvalidSecretKey:              "Некорректный секретный ключ",
	errorWrongSecretKey:                "Некорректный секретный ключ",
	errorForbiddenByRole:               "Недостаточно прав доступа для выполнения запроса",
	// analysis
	errorInvalidHighDensityCholesterol:          "Некорректное значение холестерина высокой плотности (ЛПВП)",
	errorInvalidLowDensityCholesterol:           "Некорректное значение холестерина низкой плотности (ЛПНП)",
	errorInvalidTriglycerides:                   "Некорректное значение триглицеридов",
	errorInvalidLipoprotein:                     "Некорректное значение липопротеина",
	errorInvalidHighlySensitiveCReactiveProtein: "Некорректное значение высокочувствительного С-реактивного белка",
	errorInvalidAtherogenicityCoefficient:       "Некорректное значение коэффициента атерогенности",
	errorInvalidCreatinine:                      "Некорректное значение креатинина",
	// basic indicators, score
	errorInvalidWeight:                       "Некорректное значение веса",
	errorInvalidHeight:                       "Некорректное значение роста",
	errorInvalidBodyMassIndex:                "Некорректное значение индекса массы тела (ИМТ)",
	errorInvalidWaistSize:                    "Некорректное значение объема талии",
	errorInvalidGender:                       "Некорректный пол",
	errorInvalidSBPLevel:                     "Некорректное значение уровня систолического АД",
	errorInvalidTotalCholesterolLevel:        "Некорректное значение общего холестерина",
	errorInvalidCVEventsRiskValue:            "Некорректное значение риска сердечно-сосудистых заболеваний",
	errorInvalidIdealCardiovascularAgesRange: "Некорректное значение идеального возраста сердечно-сосудистой системы",
	// auth, profile
	errorInvalidFirstName: "Некорректное значение имени",
	errorInvalidLastName:  "Некорректное значение фамилии",
	errorInvalidRegion:    "Некорректное значение региона",
	errorInvalidBirthDate: "Некорректное значение даты рождения",
	errorInvalidLogin:     "Некорректное значение логина",
	errorInvalidEmail:     "Некорректное значение электронной почты",
	errorInvalidPassword:  "Некорректное значение пароля",
	// score
	errorInvalidAge: "Некорректное значение возраста",
	// recommendations
	errorNotEnoughDataToCompileReport: "Недостаточно данных в профиле для формирования и отправки отчёта",
}

type response struct {
	Result      string `json:"result,omitempty"`
	Error       string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
	DevMessage  string `json:"devMessage,omitempty"`
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
		DevMessage:  err.Error(),
	}
}
