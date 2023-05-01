package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
)

func (r *Router) initLifestylesRoutes(customerAPI *echo.Group) {
	lifestyle := customerAPI.Group("/lifestyles", r.identifyUser, r.verifyCustomer)
	{
		lifestyle.GET("/info", r.getLifestyleInfo)
		lifestyle.PUT("/edit", r.editLifestyleInfo)
	}
}

func (r *Router) getLifestyleInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	userLifestyle, err := r.services.Lifestyle().Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, userLifestyle)
}

func (r *Router) editLifestyleInfo(c echo.Context) error {
	var reqData domain.Lifestyle
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.Lifestyle().Update(reqData); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
