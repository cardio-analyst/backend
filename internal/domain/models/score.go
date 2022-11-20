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

func (d ScoreData) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Age, validation.Min(40), validation.Max(89)),
		validation.Field(&d.Gender, validation.In(common.UserGenderMale, common.UserGenderFemale)),
		validation.Field(&d.SBPLevel, validation.Min(80.0), validation.Max(250.0)),
		validation.Field(&d.TotalCholesterolLevel, validation.Min(3.0), validation.Max(15.2)),
	)
}
