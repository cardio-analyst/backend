package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
)

const basicIndicatorsIDPathKey = "basicIndicatorsID"

func (r *Router) initBasicIndicatorsRoutes() {
	basicIndicators := r.api.Group("/basicIndicators", r.identifyUser)
	basicIndicators.GET("", r.getUserBasicIndicators)
	basicIndicators.POST("", r.createBasicIndicatorsRecord)
	basicIndicators.PUT(fmt.Sprintf("/:%v", basicIndicatorsIDPathKey), r.updateBasicIndicatorsRecord)
}

type getUserBasicIndicatorsResponse struct {
	BasicIndicators []*models.BasicIndicators `json:"basicIndicators"`
}

func (r *Router) getUserBasicIndicators(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	basicIndicators, err := r.services.BasicIndicators().FindAll(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &getUserBasicIndicatorsResponse{
		BasicIndicators: basicIndicators,
	})
}

func (r *Router) createBasicIndicatorsRecord(c echo.Context) error {
	var reqData models.BasicIndicators
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.BasicIndicators().Create(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidBasicIndicatorsData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultCreated))
}

func (r *Router) updateBasicIndicatorsRecord(c echo.Context) error {
	basicIndicatorsID, err := strconv.ParseUint(c.Param(basicIndicatorsIDPathKey), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	var reqData models.BasicIndicators
	if err = c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.ID = basicIndicatorsID
	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err = r.services.BasicIndicators().Update(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidBasicIndicatorsData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, serviceErrors.ErrBasicIndicatorsRecordNotFound):
			return c.JSON(http.StatusNotFound, newError(c, err, errorBasicIndicatorsRecordNotFound))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
