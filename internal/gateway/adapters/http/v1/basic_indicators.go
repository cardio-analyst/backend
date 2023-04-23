package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	common "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/pkg/model"
)

const basicIndicatorsIDPathKey = "basicIndicatorsID"

// possible basic indicators errors designations
const (
	errorInvalidWeight                       = "InvalidWeight"
	errorInvalidHeight                       = "InvalidHeight"
	errorInvalidBodyMassIndex                = "InvalidBodyMassIndex"
	errorInvalidWaistSize                    = "InvalidWaistSize"
	errorInvalidGender                       = "InvalidGender"
	errorInvalidSBPLevel                     = "InvalidSBPLevel"
	errorInvalidTotalCholesterolLevel        = "InvalidTotalCholesterolLevel"
	errorInvalidCVEventsRiskValue            = "InvalidCVEventsRiskValue"
	errorInvalidIdealCardiovascularAgesRange = "InvalidIdealCardiovascularAgesRange"
	errorBasicIndicatorsRecordNotFound       = "BasicIndicatorsRecordNotFound"
)

func (r *Router) initBasicIndicatorsRoutes(customerAPI *echo.Group) {
	basicIndicators := customerAPI.Group("/basicIndicators", r.identifyCustomer)
	basicIndicators.GET("", r.getUserBasicIndicators)
	basicIndicators.POST("", r.createBasicIndicatorsRecord)
	basicIndicators.PUT(fmt.Sprintf("/:%v", basicIndicatorsIDPathKey), r.updateBasicIndicatorsRecord)
}

type getUserBasicIndicatorsResponse struct {
	BasicIndicators []*common.BasicIndicators `json:"basicIndicators"`
}

func (r *Router) getUserBasicIndicators(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	basicIndicators, err := r.services.BasicIndicators().FindAll(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	criteria := model.UserCriteria{
		ID: userID,
	}

	user, err := r.services.User().GetOne(c.Request().Context(), criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	userAge := user.Age()

	for _, basicIndicator := range basicIndicators {
		if basicIndicator.CVEventsRiskValue != nil {
			riskValue := float64(*basicIndicator.CVEventsRiskValue)
			basicIndicator.Scale = r.services.Score().ResolveScale(riskValue, userAge)
		}
	}

	return c.JSON(http.StatusOK, &getUserBasicIndicatorsResponse{
		BasicIndicators: basicIndicators,
	})
}

func (r *Router) createBasicIndicatorsRecord(c echo.Context) error {
	var reqData common.BasicIndicators
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.BasicIndicators().Create(reqData); err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidWeight):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidWeight))
		case errors.Is(err, common.ErrInvalidHeight):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidHeight))
		case errors.Is(err, common.ErrInvalidBodyMassIndex):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidBodyMassIndex))
		case errors.Is(err, common.ErrInvalidWaistSize):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidWaistSize))
		case errors.Is(err, common.ErrInvalidGender):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidGender))
		case errors.Is(err, common.ErrInvalidSBPLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidSBPLevel))
		case errors.Is(err, common.ErrInvalidTotalCholesterolLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidTotalCholesterolLevel))
		case errors.Is(err, common.ErrInvalidCVEventsRiskValue):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidCVEventsRiskValue))
		case errors.Is(err, common.ErrInvalidIdealCardiovascularAgesRange):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidIdealCardiovascularAgesRange))
		case errors.Is(err, common.ErrInvalidBasicIndicatorsData):
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

	var reqData common.BasicIndicators
	if err = c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.ID = basicIndicatorsID
	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err = r.services.BasicIndicators().Update(reqData); err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidWeight):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidWeight))
		case errors.Is(err, common.ErrInvalidHeight):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidHeight))
		case errors.Is(err, common.ErrInvalidBodyMassIndex):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidBodyMassIndex))
		case errors.Is(err, common.ErrInvalidWaistSize):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidWaistSize))
		case errors.Is(err, common.ErrInvalidGender):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidGender))
		case errors.Is(err, common.ErrInvalidSBPLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidSBPLevel))
		case errors.Is(err, common.ErrInvalidTotalCholesterolLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidTotalCholesterolLevel))
		case errors.Is(err, common.ErrInvalidCVEventsRiskValue):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidCVEventsRiskValue))
		case errors.Is(err, common.ErrInvalidIdealCardiovascularAgesRange):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidIdealCardiovascularAgesRange))
		case errors.Is(err, common.ErrInvalidBasicIndicatorsData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, common.ErrBasicIndicatorsRecordNotFound):
			return c.JSON(http.StatusNotFound, newError(c, err, errorBasicIndicatorsRecordNotFound))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
