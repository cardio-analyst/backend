package v1

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

func (r *Router) initStatisticsRoutes(moderatorAPI *echo.Group) {
	moderatorAPI.GET("/statistics", r.getStatistics, r.identifyUser, r.verifyModerator)
}

type usersByRegionsItem struct {
	Region string `json:"region"`
	Value  int64  `json:"value"`
}

type getStatisticsResponse struct {
	UsersByRegions []usersByRegionsItem `json:"usersByRegions"`
}

func (r *Router) getStatistics(c echo.Context) error {
	usersByRegions, err := r.services.Statistics().UsersByRegions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	usersByRegionsItems := make([]usersByRegionsItem, 0, len(usersByRegions))
	for region, usersNum := range usersByRegions {
		usersByRegionsItems = append(usersByRegionsItems, usersByRegionsItem{
			Region: region,
			Value:  usersNum,
		})
	}

	sort.SliceStable(usersByRegionsItems, func(i, j int) bool { return usersByRegionsItems[i].Value < usersByRegionsItems[j].Value })

	return c.JSON(http.StatusOK, &getStatisticsResponse{
		UsersByRegions: usersByRegionsItems,
	})
}
