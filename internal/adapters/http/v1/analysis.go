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

const analysisIDPathKey = "analysisID"

func (r *Router) initAnalysesRoutes() {
	analyses := r.api.Group("/analyses", r.identifyUser)
	analyses.GET("", r.getUserAnalyses)
	analyses.POST("", r.createAnalysisRecord)
	analyses.PUT(fmt.Sprintf("/:%v", analysisIDPathKey), r.updateAnalysisRecord)
}

type getUserAnalysesResponse struct {
	Analyses []*models.Analysis `json:"analyses"`
}

func (r *Router) getUserAnalyses(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	analyses, err := r.services.Analysis().FindAll(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &getUserAnalysesResponse{
		Analyses: analyses,
	})
}

func (r *Router) createAnalysisRecord(c echo.Context) error {
	var reqData models.Analysis
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.Analysis().Create(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidAnalysisData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultCreated))
}

func (r *Router) updateAnalysisRecord(c echo.Context) error {
	analysisID, err := strconv.ParseUint(c.Param(analysisIDPathKey), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	var reqData models.Analysis
	if err = c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.ID = analysisID
	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err = r.services.Analysis().Update(reqData); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidAnalysisData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, serviceErrors.ErrAnalysisRecordNotFound):
			return c.JSON(http.StatusNotFound, newError(c, err, errorAnalysisRecordNotFound))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
