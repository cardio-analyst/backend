package v1

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const userIDPathKey = "userID"

const (
	queryParamRegion        = "region"
	queryParamBirthDateFrom = "birthDateFrom"
	queryParamBirthDateTo   = "birthDateTo"
)

func (r *Router) initUsersRoutes(moderatorAPI *echo.Group) {
	users := moderatorAPI.Group("/users", r.identifyUser, r.verifyModerator)
	{
		users.GET("", r.getUsers)
		users.GET(fmt.Sprintf("/:%v", userIDPathKey), r.getUser)
	}
}

type getUsersRequest struct {
	Limit int64 `query:"limit"`
	Page  int64 `query:"page"`
}

type getUsersResponseUser struct {
	ID        uint64     `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Region    string     `json:"region"`
	Login     string     `json:"login"`
	Email     string     `json:"email"`
	BirthDate model.Date `json:"birthDate,omitempty"`
}

type getUsersResponse struct {
	Users      []getUsersResponseUser `json:"users"`
	TotalPages int64                  `json:"totalPages,omitempty"`
}

func (r *Router) getUsers(c echo.Context) error {
	var reqData getUsersRequest
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	criteria := model.UserCriteria{
		Limit:  reqData.Limit,
		Page:   reqData.Page,
		Region: c.Request().URL.Query().Get(queryParamRegion),
	}

	birthDateFrom := c.Request().URL.Query().Get(queryParamBirthDateFrom)
	if birthDateFrom != "" {
		birthDateFromTime, err := time.Parse(model.DateLayout, birthDateFrom)
		if err != nil {
			return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
		}
		criteria.BirthDateFrom = model.Date{Time: birthDateFromTime}
	}

	birthDateTo := c.Request().URL.Query().Get(queryParamBirthDateTo)
	if birthDateTo != "" {
		birthDateToTime, err := time.Parse(model.DateLayout, birthDateTo)
		if err != nil {
			return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
		}
		criteria.BirthDateTo = model.Date{Time: birthDateToTime}
	}

	users, totalPages, err := r.services.User().GetList(c.Request().Context(), criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	responseUsers := make([]getUsersResponseUser, 0, len(users))
	for _, user := range users {
		responseUser := getUsersResponseUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Region:    user.Region,
			Login:     user.Login,
			Email:     user.Email,
			BirthDate: user.BirthDate,
		}

		responseUsers = append(responseUsers, responseUser)
	}

	return c.JSON(http.StatusOK, &getUsersResponse{
		Users:      responseUsers,
		TotalPages: totalPages,
	})
}

type getUserResponsePersonal struct {
	FullName  string     `json:"fullName"`
	Gender    string     `json:"gender,omitempty"`
	BirthDate model.Date `json:"birthDate,omitempty"`
	Region    string     `json:"region,omitempty"`
}

type getUserResponseHealth struct {
	Smoking               *bool    `json:"smoking,omitempty"`
	TakesStatins          bool     `json:"takesStatins,omitempty"`
	CVDPredisposed        bool     `json:"cvdPredisposed,omitempty"`
	CVEventsRiskValue     *int64   `json:"cvEventsRiskValue,omitempty"`
	SBPLevel              *float64 `json:"sbpLevel,omitempty"`
	TotalCholesterolLevel *float64 `json:"totalCholesterolLevel,omitempty"`
}

type getUserResponseLifestyle struct {
	FamilyStatus            string   `json:"familyStatus,omitempty"`
	EventsParticipation     string   `json:"eventsParticipation,omitempty"`
	PhysicalActivity        string   `json:"physicalActivity,omitempty"`
	WorkStatus              string   `json:"workStatus,omitempty"`
	AnginaScore             *int8    `json:"anginaScore,omitempty"`
	AdherenceDrugTherapy    *float64 `json:"adherenceDrugTherapy,omitempty"`
	AdherenceMedicalSupport *float64 `json:"adherenceMedicalSupport,omitempty"`
	AdherenceLifestyleMod   *float64 `json:"adherenceLifestyleMod,omitempty"`
}

type getUserResponseDiseases struct {
	HasChronicKidneyDisease bool `json:"hasChronicKidneyDisease,omitempty"`
	HasArterialHypertension bool `json:"hasArterialHypertension,omitempty"`
	HasIschemicHeartDisease bool `json:"hasIschemicHeartDisease,omitempty"`
	HasTypeTwoDiabetes      bool `json:"hasTypeTwoDiabetes,omitempty"`
	HadInfarctionOrStroke   bool `json:"hadInfarctionOrStroke,omitempty"`
	HasAtherosclerosis      bool `json:"hasAtherosclerosis,omitempty"`
	HasOtherCVD             bool `json:"hasOtherCVD,omitempty"`
}

type getUserResponseAnalyses struct {
	HighDensityCholesterol          *float64 `json:"highDensityCholesterol,omitempty"`
	LowDensityCholesterol           *float64 `json:"lowDensityCholesterol,omitempty"`
	Triglycerides                   *float64 `json:"triglycerides,omitempty"`
	Lipoprotein                     *float64 `json:"lipoprotein,omitempty"`
	HighlySensitiveCReactiveProtein *float64 `json:"highlySensitiveCReactiveProtein,omitempty"`
	AtherogenicityCoefficient       *float64 `json:"atherogenicityCoefficient,omitempty"`
	Creatinine                      *float64 `json:"creatinine,omitempty"`
	AtheroscleroticPlaquesPresence  *bool    `json:"atheroscleroticPlaquesPresence,omitempty"`
}

type getUserResponseDynamicValue struct {
	Date  model.Date `json:"date"`
	Value float64    `json:"value"`
}

type getUserResponseDynamic struct {
	Weight        []getUserResponseDynamicValue `json:"weight,omitempty"`
	Height        []getUserResponseDynamicValue `json:"height,omitempty"`
	BodyMassIndex []getUserResponseDynamicValue `json:"bodyMassIndex,omitempty"`
	WaistSize     []getUserResponseDynamicValue `json:"waistSize,omitempty"`
}

type getUserResponse struct {
	Personal  getUserResponsePersonal   `json:"personal"`
	Health    *getUserResponseHealth    `json:"health"`
	Lifestyle *getUserResponseLifestyle `json:"lifestyle"`
	Diseases  *getUserResponseDiseases  `json:"diseases"`
	Analyses  *getUserResponseAnalyses  `json:"analyses"`
	Dynamic   *getUserResponseDynamic   `json:"dynamic"`
}

func (r *Router) getUser(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param(userIDPathKey), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(c, err, errorParseRequestData))
	}

	ctx := c.Request().Context()

	g, ctx := errgroup.WithContext(ctx)
	mu := sync.Mutex{}

	resp := &getUserResponse{}

	g.Go(func() error {
		criteria := model.UserCriteria{
			ID: userID,
		}

		var user model.User
		user, err = r.services.User().GetOne(ctx, criteria)
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()

		resp.Personal.FullName = strings.TrimSpace(fmt.Sprintf("%v %v %v", user.LastName, user.FirstName, user.MiddleName))
		resp.Personal.BirthDate = user.BirthDate
		resp.Personal.Region = user.Region

		return nil
	})

	g.Go(func() error {
		var basicIndicators []*domain.BasicIndicators
		basicIndicators, err = r.services.BasicIndicators().FindAll(userID)
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()

		var (
			weightDynamic    []getUserResponseDynamicValue
			heightDynamic    []getUserResponseDynamicValue
			bmiDynamic       []getUserResponseDynamicValue
			waistSizeDynamic []getUserResponseDynamicValue
		)
		for _, basicIndicator := range basicIndicators {
			if basicIndicator.Smoking != nil {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Health.Smoking = basicIndicator.Smoking
			}
			if basicIndicator.Gender != nil {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Personal.Gender = *basicIndicator.Gender
			}
			if basicIndicator.CVEventsRiskValue != nil {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Health.CVEventsRiskValue = basicIndicator.CVEventsRiskValue
			}
			if basicIndicator.SBPLevel != nil {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Health.SBPLevel = basicIndicator.SBPLevel
			}
			if basicIndicator.TotalCholesterolLevel != nil {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Health.TotalCholesterolLevel = basicIndicator.TotalCholesterolLevel
			}

			if basicIndicator.Weight != nil {
				weightDynamic = append(weightDynamic, getUserResponseDynamicValue{
					Date:  model.Date{Time: basicIndicator.CreatedAt.Time},
					Value: *basicIndicator.Weight,
				})
			}
			if basicIndicator.Height != nil {
				heightDynamic = append(heightDynamic, getUserResponseDynamicValue{
					Date:  model.Date{Time: basicIndicator.CreatedAt.Time},
					Value: *basicIndicator.Height,
				})
			}
			if basicIndicator.BodyMassIndex != nil {
				bmiDynamic = append(bmiDynamic, getUserResponseDynamicValue{
					Date:  model.Date{Time: basicIndicator.CreatedAt.Time},
					Value: *basicIndicator.BodyMassIndex,
				})
			}
			if basicIndicator.WaistSize != nil {
				waistSizeDynamic = append(waistSizeDynamic, getUserResponseDynamicValue{
					Date:  model.Date{Time: basicIndicator.CreatedAt.Time},
					Value: *basicIndicator.WaistSize,
				})
			}
		}

		if len(weightDynamic)+len(heightDynamic)+len(bmiDynamic)+len(waistSizeDynamic) > 0 {
			resp.Dynamic = &getUserResponseDynamic{
				Weight:        weightDynamic,
				Height:        heightDynamic,
				BodyMassIndex: bmiDynamic,
				WaistSize:     waistSizeDynamic,
			}
		}

		return nil
	})

	g.Go(func() error {
		var userDiseases *domain.Diseases
		userDiseases, err = r.services.Diseases().Get(userID)
		if err != nil {
			return err
		}

		if userDiseases != nil {
			mu.Lock()
			defer mu.Unlock()

			if userDiseases.TakesStatins {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Health.TakesStatins = true
			}
			if userDiseases.CVDPredisposed {
				if resp.Health == nil {
					resp.Health = &getUserResponseHealth{}
				}
				resp.Health.CVDPredisposed = true
			}

			if userDiseases.HasChronicKidneyDisease {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HasChronicKidneyDisease = true
			}
			if userDiseases.HasArterialHypertension {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HasArterialHypertension = true
			}
			if userDiseases.HasIschemicHeartDisease {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HasIschemicHeartDisease = true
			}
			if userDiseases.HasTypeTwoDiabetes {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HasTypeTwoDiabetes = true
			}
			if userDiseases.HadInfarctionOrStroke {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HadInfarctionOrStroke = true
			}
			if userDiseases.HasAtherosclerosis {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HasAtherosclerosis = true
			}
			if userDiseases.HasOtherCVD {
				if resp.Diseases == nil {
					resp.Diseases = &getUserResponseDiseases{}
				}
				resp.Diseases.HasOtherCVD = true
			}
		}

		return nil
	})

	g.Go(func() error {
		var lifestyle *domain.Lifestyle
		lifestyle, err = r.services.Lifestyle().Get(userID)
		if err != nil {
			return err
		}

		if lifestyle != nil {
			if lifestyle.FamilyStatus != "" || lifestyle.EventsParticipation != "" || lifestyle.PhysicalActivity != "" || lifestyle.WorkStatus != "" {
				mu.Lock()
				defer mu.Unlock()

				if resp.Lifestyle == nil {
					resp.Lifestyle = &getUserResponseLifestyle{}
				}
				resp.Lifestyle.FamilyStatus = lifestyle.FamilyStatus
				resp.Lifestyle.EventsParticipation = lifestyle.EventsParticipation
				resp.Lifestyle.PhysicalActivity = lifestyle.PhysicalActivity
				resp.Lifestyle.WorkStatus = lifestyle.WorkStatus
			}
		}

		return nil
	})

	g.Go(func() error {
		var questionnaire *domain.Questionnaire
		questionnaire, err = r.services.Questionnaire().Get(userID)
		if err != nil {
			return err
		}

		if questionnaire != nil {
			mu.Lock()
			defer mu.Unlock()

			if questionnaire.AnginaScore >= 0 {
				if resp.Lifestyle == nil {
					resp.Lifestyle = &getUserResponseLifestyle{}
				}
				resp.Lifestyle.AnginaScore = &questionnaire.AnginaScore
			}
			if questionnaire.AdherenceLifestyleMod >= 0 {
				if resp.Lifestyle == nil {
					resp.Lifestyle = &getUserResponseLifestyle{}
				}
				resp.Lifestyle.AdherenceLifestyleMod = &questionnaire.AdherenceLifestyleMod
			}
			if questionnaire.AdherenceMedicalSupport >= 0 {
				if resp.Lifestyle == nil {
					resp.Lifestyle = &getUserResponseLifestyle{}
				}
				resp.Lifestyle.AdherenceMedicalSupport = &questionnaire.AdherenceMedicalSupport
			}
			if questionnaire.AdherenceDrugTherapy >= 0 {
				if resp.Lifestyle == nil {
					resp.Lifestyle = &getUserResponseLifestyle{}
				}
				resp.Lifestyle.AdherenceDrugTherapy = &questionnaire.AdherenceDrugTherapy
			}
		}

		return nil
	})

	g.Go(func() error {
		var analyses []*domain.Analysis
		analyses, err = r.services.Analysis().FindAll(userID)
		if err != nil {
			return err
		}

		if len(analyses) > 0 {
			mu.Lock()
			defer mu.Unlock()

			for _, analysis := range analyses {
				if analysis.HighDensityCholesterol != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.HighDensityCholesterol = analysis.HighDensityCholesterol
				}
				if analysis.LowDensityCholesterol != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.LowDensityCholesterol = analysis.LowDensityCholesterol
				}
				if analysis.Triglycerides != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.Triglycerides = analysis.Triglycerides
				}
				if analysis.Lipoprotein != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.Lipoprotein = analysis.Lipoprotein
				}
				if analysis.HighlySensitiveCReactiveProtein != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.HighlySensitiveCReactiveProtein = analysis.HighlySensitiveCReactiveProtein
				}
				if analysis.AtherogenicityCoefficient != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.AtherogenicityCoefficient = analysis.AtherogenicityCoefficient
				}
				if analysis.Creatinine != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.Creatinine = analysis.Creatinine
				}
				if analysis.AtheroscleroticPlaquesPresence != nil {
					if resp.Analyses == nil {
						resp.Analyses = &getUserResponseAnalyses{}
					}
					resp.Analyses.AtheroscleroticPlaquesPresence = analysis.AtheroscleroticPlaquesPresence
				}
			}
		}

		return nil
	})

	if err = g.Wait(); err != nil {
		return c.JSON(http.StatusInternalServerError, newError(c, err, errorInternal))
	}

	return c.JSON(http.StatusOK, resp)
}
