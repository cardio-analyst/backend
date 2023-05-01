package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
)

func (r *Router) initQuestionnaireRoutes(customerAPI *echo.Group) {
	tests := customerAPI.Group("/tests", r.identifyUser, r.verifyCustomer)
	{
		angina := tests.Group("/angina-rose")
		{
			angina.GET("/info", r.anginaRoseInfo)
			angina.PUT("/edit", r.anginaRoseEdit)
		}

		adherence := tests.Group("/treatment-adherence")
		{
			adherence.GET("/info", r.treatmentAdherenceInfo)
			adherence.PUT("/edit", r.treatmentAdherenceEdit)
		}
	}
}

type anginaRoseInfoResponse struct {
	AnginaScore int8 `json:"anginaScore" db:"angina_score"`
}

func (r *Router) anginaRoseInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	questionnaire, err := r.services.Questionnaire().Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &anginaRoseInfoResponse{
		AnginaScore: questionnaire.AnginaScore,
	})
}

func (r *Router) anginaRoseEdit(c echo.Context) error {
	var reqData domain.Questionnaire
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.Questionnaire().UpdateAnginaRose(reqData); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}

type treatmentAdherenceInfoResponse struct {
	AdherenceDrugTherapy    float64 `json:"adherenceDrugTherapy"`
	AdherenceMedicalSupport float64 `json:"adherenceMedicalSupport"`
	AdherenceLifestyleMod   float64 `json:"adherenceLifestyleMod"`
}

func (r *Router) treatmentAdherenceInfo(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)

	questionnaire, err := r.services.Questionnaire().Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &treatmentAdherenceInfoResponse{
		AdherenceDrugTherapy:    questionnaire.AdherenceDrugTherapy,
		AdherenceMedicalSupport: questionnaire.AdherenceMedicalSupport,
		AdherenceLifestyleMod:   questionnaire.AdherenceLifestyleMod,
	})
}

func (r *Router) treatmentAdherenceEdit(c echo.Context) error {
	var reqData domain.Questionnaire
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	reqData.UserID = c.Get(ctxKeyUserID).(uint64)

	if err := r.services.Questionnaire().UpdateTreatmentAdherence(reqData); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
