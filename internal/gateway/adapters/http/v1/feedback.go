package v1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (r *Router) initFeedbackRoutes() {
	feedback := r.api.Group(fmt.Sprintf("/:%v/feedback", userRolePathKey), r.identifyUser, r.parseUserRole)
	{
		feedback.GET("", r.getFeedbacks, r.verifyModerator)
		feedback.POST("/send", r.sendFeedback, r.verifyCustomer)
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
}

func (r *Router) sendFeedback(c echo.Context) error {
	var reqData sendFeedbackRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	userID := c.Get(ctxKeyUserID).(uint64)

	user, err := r.services.User().GetOne(c.Request().Context(), model.UserCriteria{ID: userID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	if err = r.services.Feedback().Send(reqData.Mark, reqData.Message, user); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, newResult(resultSent))
}
