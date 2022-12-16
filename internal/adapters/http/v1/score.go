package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serviceErrors "github.com/cardio-analyst/backend/internal/domain/errors"
	"github.com/cardio-analyst/backend/internal/domain/models"
)

// possible score errors designations
const (
	errorInvalidAge           = "InvalidAge"
	errorNotEnoughInformation = "NotEnoughInformation"
)

func (r *Router) initScoreRoutes() {
	score := r.api.Group("/score", r.identifyUser)
	score.GET("/cveRisk", r.cveRisk)
	score.GET("/idealAge", r.idealAge)
}

type getCVERiskResponse struct {
	Value uint64 `json:"value"`
	Scale string `json:"scale"`
}

func (r *Router) cveRisk(c echo.Context) error {
	var reqData models.ScoreData
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	userID := c.Get(ctxKeyUserID).(uint64)

	criteria := models.UserCriteria{
		ID: &userID,
	}

	user, err := r.services.User().Get(criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	reqData.Age = user.Age()

	riskValue, scale, err := r.services.Score().GetCVERisk(reqData)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidAge):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidAge))
		case errors.Is(err, serviceErrors.ErrInvalidGender):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidGender))
		case errors.Is(err, serviceErrors.ErrInvalidSBPLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidSBPLevel))
		case errors.Is(err, serviceErrors.ErrInvalidTotalCholesterolLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidTotalCholesterolLevel))
		case errors.Is(err, serviceErrors.ErrInvalidScoreData):
			return c.JSON(http.StatusUnprocessableEntity, newError(c, err, errorNotEnoughInformation))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, &getCVERiskResponse{
		Value: riskValue,
		Scale: scale,
	})
}

type getIdealAgeResponse struct {
	Value string `json:"value"`
	Scale string `json:"scale"`
}

func (r *Router) idealAge(c echo.Context) error {
	var reqData models.ScoreData
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	userID := c.Get(ctxKeyUserID).(uint64)

	criteria := models.UserCriteria{
		ID: &userID,
	}

	user, err := r.services.User().Get(criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	reqData.Age = user.Age()

	agesRange, scale, err := r.services.Score().GetIdealAge(reqData)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidAge):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidAge))
		case errors.Is(err, serviceErrors.ErrInvalidGender):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidGender))
		case errors.Is(err, serviceErrors.ErrInvalidSBPLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidSBPLevel))
		case errors.Is(err, serviceErrors.ErrInvalidTotalCholesterolLevel):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidTotalCholesterolLevel))
		case errors.Is(err, serviceErrors.ErrInvalidScoreData):
			return c.JSON(http.StatusUnprocessableEntity, newError(c, err, errorNotEnoughInformation))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, &getIdealAgeResponse{
		Value: agesRange,
		Scale: scale,
	})
}
