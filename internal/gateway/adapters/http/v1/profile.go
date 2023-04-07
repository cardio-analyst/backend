package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/pkg/model"
)

func (r *Router) initProfileRoutes() {
	profile := r.api.Group("/profile", r.identifyUser)
	profile.GET("/info", r.getProfileInfo)
	profile.PUT("/edit", r.editProfileInfo)
}

func (r *Router) getProfileInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	criteria := model.UserCriteria{
		ID: userID,
	}

	user, err := r.services.User().GetOne(c.Request().Context(), criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, user)
}

func (r *Router) editProfileInfo(c echo.Context) error {
	var reqData model.User
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.ID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.User().Update(c.Request().Context(), reqData); err != nil {
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
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
