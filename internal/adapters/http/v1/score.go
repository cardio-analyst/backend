package v1

import (
	"github.com/cardio-analyst/backend/internal/domain/common"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func (r *Router) initScoreRoutes() {
	cveRisk := r.api.Group("/score", r.identifyUser)
	cveRisk.GET("/cveRisk", r.getCveRisk)
}

func (r *Router) getCveRisk(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	criteria := models.UserCriteria{
		ID: &userID,
	}

	user, err := r.services.User().Get(criteria)

	if err != nil {
		return err
	}

	// TODO added validate data
	genderParam := c.QueryParam("gender")
	smokingParam := c.QueryParam("smoking")
	sbpLevelParam := c.QueryParam("sbpLevel")
	totalCholesterolLevelParam := c.QueryParam("totalCholesterolLevel")

	smoking, err := strconv.ParseBool(smokingParam)

	if err != nil {
		return err
	}

	sbpLevel, err := strconv.ParseUint(sbpLevelParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	totalCholesterolLevel, err := strconv.ParseFloat(totalCholesterolLevelParam, 64)

	age := common.GetAge(user.BirthDate.Time, time.Now())

	cveRiskData := models.CveRiskData{
		Age:                   age,
		Gender:                genderParam,
		Smoking:               smoking,
		SbpLevel:              sbpLevel,
		TotalCholesterolLevel: totalCholesterolLevel,
	}

	riskValue, err := r.services.Score().GetCveRisk(cveRiskData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, riskValue)
}
