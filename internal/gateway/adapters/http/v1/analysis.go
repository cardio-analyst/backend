package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
)

const analysisIDPathKey = "analysisID"

// possible analysis errors designations
const (
	errorAnalysisRecordNotFound                 = "AnalysisRecordNotFound"
	errorInvalidHighDensityCholesterol          = "InvalidHighDensityCholesterol"
	errorInvalidLowDensityCholesterol           = "InvalidLowDensityCholesterol"
	errorInvalidTriglycerides                   = "InvalidTriglycerides"
	errorInvalidLipoprotein                     = "InvalidLipoprotein"
	errorInvalidHighlySensitiveCReactiveProtein = "InvalidHighlySensitiveCReactiveProtein"
	errorInvalidAtherogenicityCoefficient       = "InvalidAtherogenicityCoefficient"
	errorInvalidCreatinine                      = "InvalidCreatinine"
)

func (r *Router) initAnalysesRoutes(customerAPI *echo.Group) {
	analyses := customerAPI.Group("/analyses", r.identifyCustomer)
	analyses.GET("", r.getUserAnalyses)
	analyses.POST("", r.createAnalysisRecord)
	analyses.PUT(fmt.Sprintf("/:%v", analysisIDPathKey), r.updateAnalysisRecord)
}

type getUserAnalysesResponse struct {
	Analyses []*domain.Analysis `json:"analyses"`
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
	var reqData domain.Analysis
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.Analysis().Create(reqData); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidHighDensityCholesterol):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidHighDensityCholesterol))
		case errors.Is(err, domain.ErrInvalidLowDensityCholesterol):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLowDensityCholesterol))
		case errors.Is(err, domain.ErrInvalidTriglycerides):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidTriglycerides))
		case errors.Is(err, domain.ErrInvalidLipoprotein):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLipoprotein))
		case errors.Is(err, domain.ErrInvalidHighlySensitiveCReactiveProtein):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidHighlySensitiveCReactiveProtein))
		case errors.Is(err, domain.ErrInvalidAtherogenicityCoefficient):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidAtherogenicityCoefficient))
		case errors.Is(err, domain.ErrInvalidCreatinine):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidCreatinine))
		case errors.Is(err, domain.ErrInvalidAnalysisData):
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

	var reqData domain.Analysis
	if err = c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.ID = analysisID
	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err = r.services.Analysis().Update(reqData); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidHighDensityCholesterol):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidHighDensityCholesterol))
		case errors.Is(err, domain.ErrInvalidLowDensityCholesterol):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLowDensityCholesterol))
		case errors.Is(err, domain.ErrInvalidTriglycerides):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidTriglycerides))
		case errors.Is(err, domain.ErrInvalidLipoprotein):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidLipoprotein))
		case errors.Is(err, domain.ErrInvalidHighlySensitiveCReactiveProtein):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidHighlySensitiveCReactiveProtein))
		case errors.Is(err, domain.ErrInvalidAtherogenicityCoefficient):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidAtherogenicityCoefficient))
		case errors.Is(err, domain.ErrInvalidCreatinine):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidCreatinine))
		case errors.Is(err, domain.ErrInvalidAnalysisData):
			return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
		case errors.Is(err, domain.ErrAnalysisRecordNotFound):
			return c.JSON(http.StatusNotFound, newError(c, err, errorAnalysisRecordNotFound))
		default:
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
