package v1

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const feedbackIDPathKey = "feedbackID"

const errorFeedbackNotFound = "FeedbackNotFound"

const (
	orderingDisabled = ""
	orderingASC      = "ASC"
	orderingDESC     = "DESC"
)

const feedbackCreatedAtLayout = "2006.01.02 15:04:05"

func (r *Router) initFeedbackRoutes() {
	feedback := r.api.Group(fmt.Sprintf("/:%v/feedback", userRolePathKey), r.identifyUser, r.parseUserRole)
	{
		feedback.POST("/send", r.sendFeedback, r.verifyCustomer)

		feedback.GET("", r.getFeedbacks, r.verifyModerator)
		feedback.PUT(fmt.Sprintf("/:%v", feedbackIDPathKey), r.toggleFeedbackViewed, r.verifyModerator)
	}
}

type getFeedbacksRequest struct {
	Limit           int64  `query:"limit"`
	Page            int64  `query:"page"`
	Viewed          *bool  `query:"viewed"`
	MarkOrdering    string `query:"markOrdering"`
	VersionOrdering string `query:"versionOrdering"`
}

type getFeedbacksResponseItem struct {
	ID             uint64 `json:"id"`
	UserID         uint64 `json:"userId"`
	UserFirstName  string `json:"userFirstName"`
	UserLastName   string `json:"userLastName"`
	UserMiddleName string `json:"userMiddleName,omitempty"`
	UserLogin      string `json:"userLogin"`
	UserEmail      string `json:"userEmail"`
	Mark           int16  `json:"mark"`
	Message        string `json:"message,omitempty"`
	Version        string `json:"version"`
	Viewed         bool   `json:"viewed"`
	CreatedAt      string `json:"createdAt"`
}

type getFeedbacksResponse struct {
	Feedbacks  []getFeedbacksResponseItem `json:"feedbacks"`
	TotalPages int64                      `json:"totalPages,omitempty"`
}

func (r *Router) getFeedbacks(c echo.Context) error {
	var reqData getFeedbacksRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	criteria := model.FeedbackCriteria{
		Limit:  reqData.Limit,
		Page:   reqData.Page,
		Viewed: reqData.Viewed,
	}

	switch reqData.MarkOrdering {
	case orderingDisabled:
		criteria.MarkOrdering = model.OrderingTypeDisabled
	case orderingASC:
		criteria.MarkOrdering = model.OrderingTypeASC
	case orderingDESC:
		criteria.MarkOrdering = model.OrderingTypeDESC
	default:
		err := fmt.Errorf("undefined mark ordering type: %q", reqData.MarkOrdering)
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	switch reqData.VersionOrdering {
	case orderingDisabled:
		criteria.VersionOrdering = model.OrderingTypeDisabled
	case orderingASC:
		criteria.VersionOrdering = model.OrderingTypeASC
	case orderingDESC:
		criteria.VersionOrdering = model.OrderingTypeDESC
	default:
		err := fmt.Errorf("undefined version ordering type: %q", reqData.VersionOrdering)
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	feedbacks, totalPages, err := r.services.Feedback().FindAll(criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	responseItems := make([]getFeedbacksResponseItem, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		responseItems = append(responseItems, getFeedbacksResponseItem{
			ID:             feedback.ID,
			UserID:         feedback.UserID,
			UserFirstName:  feedback.UserFirstName,
			UserLastName:   feedback.UserLastName,
			UserMiddleName: feedback.UserMiddleName,
			UserLogin:      feedback.UserLogin,
			UserEmail:      feedback.UserEmail,
			Mark:           feedback.Mark,
			Message:        feedback.Message,
			Version:        feedback.Version,
			Viewed:         feedback.Viewed,
			CreatedAt:      feedback.CreatedAt.Time.Format(feedbackCreatedAtLayout),
		})
	}

	return c.JSON(http.StatusOK, &getFeedbacksResponse{
		Feedbacks:  responseItems,
		TotalPages: totalPages,
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
		validation.Field(&r.Version, validation.Match(regexp.MustCompile("^([0-9]+.){2}[0-9]+$"))),
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
