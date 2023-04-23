package v1

import (
	"errors"
	"net/http"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/pkg/model"
)

// possible recommendations errors designations
const (
	errorNotEnoughDataToCompileReport = "NotEnoughDataToCompileReport"
)

var errNoOneToSendReport = errors.New("there is no one to send report to")

func (r *Router) initRecommendationsRoutes(customerAPI *echo.Group) {
	recommendations := customerAPI.Group("/recommendations", r.identifyCustomer)
	recommendations.GET("", r.getRecommendations)
	recommendations.POST("/send", r.sendRecommendations)
}

type getRecommendationsResponse struct {
	Recommendations []*domain.Recommendation `json:"recommendations"`
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

	user, err := r.services.User().GetOne(c.Request().Context(), model.UserCriteria{ID: userID})
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

	reportFilePath, err := r.services.Report(domain.PDF).GenerateReport(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotEnoughDataToCompileReport) {
			return c.JSON(http.StatusBadRequest, newError(c, err, errorNotEnoughDataToCompileReport))
		}
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}
	defer func() {
		if err = os.Remove(reportFilePath); err != nil {
			log.Warnf("removing report file %q: %v", reportFilePath, err)
		}
	}()

	if err = r.services.Email().SendReport(receivers, reportFilePath, user); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultEmailSent))
}
