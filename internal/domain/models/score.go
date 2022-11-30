package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/cardio-analyst/backend/internal/domain/common"
)

type ScoreData struct {
	Age                   int     // receive from user data
	Gender                string  `query:"gender"`
	Smoking               bool    `query:"smoking"`
	SBPLevel              float64 `query:"sbpLevel"`
	TotalCholesterolLevel float64 `query:"totalCholesterolLevel"`
}

func ExtractScoreDataFrom(indicators []*BasicIndicators) ScoreData {
	var data ScoreData
	for _, basicIndicator := range indicators {
		if basicIndicator.Smoking != nil && *basicIndicator.Smoking {
			data.Smoking = true
		}
		if basicIndicator.Gender != nil && data.Gender == "" {
			data.Gender = *basicIndicator.Gender
		}
		if basicIndicator.SBPLevel != nil && data.SBPLevel == 0 {
			data.SBPLevel = *basicIndicator.SBPLevel
		}
		if basicIndicator.TotalCholesterolLevel != nil && data.TotalCholesterolLevel == 0 {
			data.TotalCholesterolLevel = *basicIndicator.TotalCholesterolLevel
		}

		// fastest break condition
		if data.Gender != "" && data.SBPLevel != 0 && data.TotalCholesterolLevel != 0 {
			break
		}
	}
	return data
}

func (d ScoreData) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Age, validation.Min(40), validation.Max(89)),
		validation.Field(&d.Gender, validation.In(common.UserGenderMale, common.UserGenderFemale)),
		validation.Field(&d.SBPLevel, validation.Min(100.0), validation.Max(179.0)),
		validation.Field(&d.TotalCholesterolLevel, validation.Min(3.0), validation.Max(6.9)),
	)
}
