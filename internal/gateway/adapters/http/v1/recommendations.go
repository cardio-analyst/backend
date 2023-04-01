package v1

import (
	"errors"
	models2 "github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"net/http"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

var errNoOneToSendReport = errors.New("there is no one to send report to")

func (r *Router) initRecommendationsRoutes() {
	recommendations := r.api.Group("/recommendations", r.identifyUser)
	recommendations.GET("", r.getRecommendations)
	recommendations.POST("/send", r.sendRecommendations)
}

type getRecommendationsResponse struct {
	Recommendations []*models2.Recommendation `json:"recommendations"`
}

func (r *Router) getRecommendations(c echo.Context) error {
	userID := c.Get(ctxKeyUserID).(uint64)
	recommendations, err := r.services.Recommendations().GetRecommendations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &getRecommendationsResponse{
		Recommendations: recommendations,
	})
}

type sendRecommendationsRequest struct {
	Receiver   string `json:"receiver"`
	SendMyself bool   `json:"sendMyself"`
}

func (r sendRecommendationsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Receiver, validation.When(
			r.Receiver != "",
			is.Email,
		)),
	)
}

func (r *Router) sendRecommendations(c echo.Context) error {
	var reqData sendRecommendationsRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	if reqData.Receiver == "" && !reqData.SendMyself {
		return c.JSON(http.StatusUnprocessableEntity, newError(c, errNoOneToSendReport, errorNotEnoughInformation))
	}

	if err := reqData.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidEmail))
	}

	userID := c.Get(ctxKeyUserID).(uint64)

	user, err := r.services.User().Get(models2.UserCriteria{ID: &userID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	receivers := make([]string, 0, 2)
	if reqData.Receiver != "" {
		receivers = append(receivers, reqData.Receiver)
	}
	if reqData.SendMyself && reqData.Receiver != user.Email {
		receivers = append(receivers, user.Email)
	}

	reportPath, err := r.services.Report(models2.PDF).GenerateReport(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}
	defer os.Remove(reportPath)

	if err = r.services.Email().SendReport(receivers, reportPath, *user); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultEmailSent))
}
