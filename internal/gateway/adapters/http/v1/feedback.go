package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const feedbackIDPathKey = "feedbackID"

const errorFeedbackNotFound = "FeedbackNotFound"

func (r *Router) initFeedbackRoutes() {
	feedback := r.api.Group(fmt.Sprintf("/:%v/feedback", userRolePathKey), r.identifyUser, r.parseUserRole)
	{
		feedback.POST("/send", r.sendFeedback, r.verifyCustomer)

		feedback.GET("", r.getFeedbacks, r.verifyModerator)
		feedback.PUT(fmt.Sprintf("/:%v", feedbackIDPathKey), r.toggleFeedbackViewed, r.verifyModerator)
	}
}

type getFeedbacksResponse struct {
	Feedbacks []model.Feedback `json:"feedbacks"`
}

func (r *Router) getFeedbacks(c echo.Context) error {
	feedbacks, err := r.services.Feedback().FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, &getFeedbacksResponse{
		Feedbacks: feedbacks,
	})
}

type sendFeedbackRequest struct {
	Mark    int16  `json:"mark"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func (r sendFeedbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Mark, validation.Min(0), validation.Max(5)),
		validation.Field(&r.Version, validation.Required),
	)
}

func (r *Router) sendFeedback(c echo.Context) error {
	var reqData sendFeedbackRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}
	if err := reqData.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
	}

	userID := c.Get(ctxKeyUserID).(uint64)

	user, err := r.services.User().GetOne(c.Request().Context(), model.UserCriteria{ID: userID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	if err = r.services.Feedback().Send(reqData.Mark, reqData.Message, reqData.Version, user); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultSent))
}

func (r *Router) toggleFeedbackViewed(c echo.Context) error {
	feedbackID, err := strconv.ParseUint(c.Param(feedbackIDPathKey), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorInvalidRequestData))
	}

	if err = r.services.Feedback().ToggleFeedbackViewed(feedbackID); err != nil {
		if errors.Is(err, model.ErrFeedbackNotFound) {
			return c.JSON(http.StatusBadRequest, newError(c, err, errorFeedbackNotFound))
		}
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultUpdated))
}
