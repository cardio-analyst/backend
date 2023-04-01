package v1

import (
	"errors"
	errors2 "github.com/cardio-analyst/backend/internal/gateway/domain/errors"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Router) initProfileRoutes() {
	profile := r.api.Group("/profile", r.identifyUser)
	profile.GET("/info", r.getProfileInfo)
	profile.PUT("/edit", r.editProfileInfo)
}

func (r *Router) getProfileInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	criteria := models.UserCriteria{
		ID: &userID,
	}

	user, err := r.services.User().Get(criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, user)
}

func (r *Router) editProfileInfo(c echo.Context) error {
	var reqData models.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.ID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.User().Update(reqData); err != nil {
		switch {
		case errors.Is(err, errors2.ErrInvalidFirstName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidFirstName))
		case errors.Is(err, errors2.ErrInvalidLastName):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLastName))
		case errors.Is(err, errors2.ErrInvalidRegion):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRegion))
		case errors.Is(err, errors2.ErrInvalidBirthDate):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidBirthDate))
		case errors.Is(err, errors2.ErrInvalidLogin):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLogin))
		case errors.Is(err, errors2.ErrInvalidEmail):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidEmail))
		case errors.Is(err, errors2.ErrInvalidPassword):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidPassword))
		case errors.Is(err, errors2.ErrInvalidUserData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, errors2.ErrUserLoginAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorLoginAlreadyOccupied))
		case errors.Is(err, errors2.ErrUserEmailAlreadyOccupied):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorEmailAlreadyOccupied))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
