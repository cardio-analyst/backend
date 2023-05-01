package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
)

func (r *Router) initDiseasesRoutes(customerAPI *echo.Group) {
	diseases := customerAPI.Group("/diseases", r.identifyUser, r.verifyCustomer)
	{
		diseases.GET("/info", r.getDiseasesInfo)
		diseases.PUT("/edit", r.editDiseasesInfo)
	}
}

func (r *Router) getDiseasesInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	userDiseases, err := r.services.Diseases().Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, userDiseases)
}

func (r *Router) editDiseasesInfo(c echo.Context) error {
	var reqData domain.Diseases
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.Diseases().Update(reqData); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
