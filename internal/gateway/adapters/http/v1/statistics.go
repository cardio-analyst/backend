package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

func (r *Router) initStatisticsRoutes(moderatorAPI *echo.Group) {
	moderatorAPI.GET("/statistics", r.getStatistics, r.identifyUser, r.verifyModerator)
}

type usersByRegionsItem struct {
	Region string `json:"region"`
	Value  int64  `json:"value"`
}

type diseasesItem struct {
	Disease string `json:"disease"`
	Value   int64  `json:"value"`
}

type sbpItem struct {
	Range string `json:"range"`
	Value int64  `json:"value"`
}

type idealCardiovascularAgesRangeItem struct {
	Range string `json:"range"`
	Value int64  `json:"value"`
}

type getStatisticsResponse struct {
	UsersByRegions                  []usersByRegionsItem               `json:"usersByRegions"`
	DiseasesByUsers                 []diseasesItem                     `json:"diseasesByUsers"`
	SBPByUsers                      []sbpItem                          `json:"sbpByUsers"`
	CardiovascularAgesRangesByUsers []idealCardiovascularAgesRangeItem `json:"cardiovascularAgesRangesByUsers"`
}

func (r *Router) getStatistics(c echo.Context) error {
	region := c.Request().URL.Query().Get(queryParamRegion)

	var (
		usersByRegionsItems []usersByRegionsItem
		regionUsers         map[uint64]bool
	)
	if region != "" {
		criteria := model.UserCriteria{
			Region: region,
		}

		users, _, err := r.services.User().GetList(c.Request().Context(), criteria)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}

		regionUsers = make(map[uint64]bool, len(users))
		for _, user := range users {
			regionUsers[user.ID] = true
		}
	} else {
		usersByRegions, err := r.services.Statistics().UsersByRegions(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
		}

		usersByRegionsItems = make([]usersByRegionsItem, 0, len(usersByRegions))
		for usersRegion, usersNum := range usersByRegions {
			if usersNum == 0 {
				continue
			}
			usersByRegionsItems = append(usersByRegionsItems, usersByRegionsItem{
				Region: usersRegion,
				Value:  usersNum,
			})
		}
		if len(usersByRegionsItems) == 0 {
			usersByRegionsItems = nil
		}
	}

	diseasesByUsers, err := r.services.Statistics().DiseasesByUsers(region, regionUsers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	diseasesByUsersItems := make([]diseasesItem, 0, len(diseasesByUsers))
	for disease, usersNum := range diseasesByUsers {
		if usersNum == 0 {
			continue
		}
		diseasesByUsersItems = append(diseasesByUsersItems, diseasesItem{
			Disease: disease,
			Value:   usersNum,
		})
	}
	if len(diseasesByUsersItems) == 0 {
		diseasesByUsersItems = nil
	}

	sbpByUsers, err := r.services.Statistics().SBPByUsers(region, regionUsers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	sbpByUsersItems := make([]sbpItem, 0, len(sbpByUsers))
	for sbpRange, usersNum := range sbpByUsers {
		if usersNum == 0 {
			continue
		}
		sbpByUsersItems = append(sbpByUsersItems, sbpItem{
			Range: sbpRange,
			Value: usersNum,
		})
	}
	if len(sbpByUsersItems) == 0 {
		sbpByUsersItems = nil
	}

	cardiovascularAgesRangesByUsers, err := r.services.Statistics().IdealCardiovascularAgesRangesByUsers(region, regionUsers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	cardiovascularAgesRangesByUsersItems := make([]idealCardiovascularAgesRangeItem, 0, len(cardiovascularAgesRangesByUsers))
	for cardiovascularAgesRange, usersNum := range cardiovascularAgesRangesByUsers {
		if usersNum == 0 {
			continue
		}
		cardiovascularAgesRangesByUsersItems = append(cardiovascularAgesRangesByUsersItems, idealCardiovascularAgesRangeItem{
			Range: cardiovascularAgesRange,
			Value: usersNum,
		})
	}
	if len(cardiovascularAgesRangesByUsersItems) == 0 {
		cardiovascularAgesRangesByUsersItems = nil
	}

	return c.JSON(http.StatusOK, &getStatisticsResponse{
		UsersByRegions:                  usersByRegionsItems,
		DiseasesByUsers:                 diseasesByUsersItems,
		SBPByUsers:                      sbpByUsersItems,
		CardiovascularAgesRangesByUsers: cardiovascularAgesRangesByUsersItems,
	})
}
