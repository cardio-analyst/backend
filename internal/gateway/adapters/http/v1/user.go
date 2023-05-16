package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (r *Router) initUsersRoutes(moderatorAPI *echo.Group) {
	users := moderatorAPI.Group("/users", r.identifyUser, r.verifyModerator)
	{
		users.GET("", r.getUsers)
	}
}

type getUsersRequest struct {
	Limit int64 `query:"limit"`
	Page  int64 `query:"page"`
}

type getUsersResponseUser struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Region    string `json:"region"`
	Login     string `json:"login"`
	Email     string `json:"email"`
}

type getUsersResponse struct {
	Users       []getUsersResponseUser `json:"users"`
	HasNextPage bool                   `json:"hasNextPage,omitempty"`
}

func (r *Router) getUsers(c echo.Context) error {
	var reqData getUsersRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	criteria := model.UserCriteria{
		Limit: reqData.Limit,
		Page:  reqData.Page,
	}

	users, hasNextPage, err := r.services.User().GetList(c.Request().Context(), criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	responseUsers := make([]getUsersResponseUser, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, getUsersResponseUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Region:    user.Region,
			Login:     user.Login,
			Email:     user.Email,
		})
	}

	return c.JSON(http.StatusOK, &getUsersResponse{
		Users:       responseUsers,
		HasNextPage: hasNextPage,
	})
}
