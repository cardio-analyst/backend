package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

// possible auth errors designations
const (
	errorInvalidFirstName    = "InvalidFirstName"
	errorInvalidLastName     = "InvalidLastName"
	errorInvalidRegion       = "InvalidRegion"
	errorInvalidBirthDate    = "InvalidBirthDate"
	errorInvalidLogin        = "InvalidLogin"
	errorInvalidEmail        = "InvalidEmail"
	errorInvalidPassword     = "InvalidPassword"
	errorRefreshTokenExpired = "RefreshTokenExpired"
	errorWrongRefreshToken   = "WrongRefreshToken"
	errorWrongCredentials    = "WrongCredentials"
	errorIPNotAllowed        = "IPNotAllowed"
	errorInvalidSecretKey    = "InvalidSecretKey"
	errorWrongSecretKey      = "WrongSecretKey"
)

func (r *Router) initAuthRoutes() {
	auth := r.api.Group(fmt.Sprintf("/:%v/auth", userRolePathKey), r.parseUserRole)
	{
		auth.POST("/signUp", r.signUp)
		auth.POST("/signIn", r.signIn)
		auth.POST("/refreshTokens", r.refreshTokens)
	}
}

func (r *Router) signUp(c echo.Context) error {
	var reqData model.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	userRole := c.Get(userRolePathKey).(model.UserRole)
	if userRole == model.UserRoleAdministrator {
		err := errors.New("forbidden to create the administrator")
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}
	reqData.Role = userRole

	if err := r.services.Auth().RegisterUser(c.Request().Context(), reqData); err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidFirstName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidFirstName))
		case errors.Is(err, model.ErrInvalidLastName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLastName))
		case errors.Is(err, model.ErrInvalidRegion):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRegion))
		case errors.Is(err, model.ErrInvalidBirthDate):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidBirthDate))
		case errors.Is(err, model.ErrInvalidLogin):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLogin))
		case errors.Is(err, model.ErrInvalidEmail):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidEmail))
		case errors.Is(err, model.ErrInvalidPassword):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidPassword))
		case errors.Is(err, model.ErrInvalidUserData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, model.ErrUserLoginAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorLoginAlreadyOccupied))
		case errors.Is(err, model.ErrUserEmailAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorEmailAlreadyOccupied))
		case errors.Is(err, model.ErrInvalidSecretKey):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidSecretKey))
		case errors.Is(err, model.ErrWrongSecretKey):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongSecretKey))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultRegistered))
}

type signInRequest struct {
	LoginOrEmail string `json:"loginOrEmail"`
	Password     string `json:"password"`
}

func (r *Router) signIn(c echo.Context) error {
	var reqData signInRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	credentials := model.Credentials{
		Password: reqData.Password,
	}

	_, err := mail.ParseAddress(reqData.LoginOrEmail)
	if err == nil {
		credentials.Email = reqData.LoginOrEmail
	} else {
		credentials.Login = reqData.LoginOrEmail
	}

	userRole := c.Get(userRolePathKey).(model.UserRole)

	tokens, err := r.services.Auth().GetTokens(c.Request().Context(), credentials, c.RealIP(), userRole)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCredentials):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, model.ErrWrongCredentials):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongCredentials))
		case errors.Is(err, model.ErrForbiddenByRole):
			return c.JSON(http.StatusForbidden, newError(c, err, errorForbiddenByRole))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

type refreshTokensRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r *Router) refreshTokens(c echo.Context) error {
	var reqData refreshTokensRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	userRole := c.Get(userRolePathKey).(model.UserRole)

	tokens, err := r.services.Auth().RefreshTokens(c.Request().Context(), reqData.RefreshToken, c.RealIP(), userRole)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrWrongToken):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorWrongRefreshToken))
		case errors.Is(err, model.ErrTokenIsExpired):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorRefreshTokenExpired))
		case errors.Is(err, model.ErrIPIsNotInWhitelist):
			return c.JSON(http.StatusForbidden, newError(c, err, errorIPNotAllowed))
		case errors.Is(err, model.ErrForbiddenByRole):
			return c.JSON(http.StatusForbidden, newError(c, err, errorForbiddenByRole))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

type generateSecretKeyRequest struct {
	UserLogin string `json:"userLogin"`
	UserEmail string `json:"userEmail"`
}

type generateSecretKeyResponse struct {
	SecretKey string `json:"secretKey"`
}

func (r *Router) generateSecretKey(c echo.Context) error {
	var reqData generateSecretKeyRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	secretKey, err := r.services.Auth().GenerateSecretKey(c.Request().Context(), reqData.UserLogin, reqData.UserEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	resp := &generateSecretKeyResponse{
		SecretKey: secretKey,
	}

	return c.JSON(http.StatusOK, resp)
}
