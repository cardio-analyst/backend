package models

import (
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/cardio-analyst/backend/internal/domain/common"
	"github.com/cardio-analyst/backend/internal/domain/errors"
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
	err := validation.ValidateStruct(&d,
		validation.Field(&d.Age, validation.Min(40), validation.Max(89)),
		validation.Field(&d.Gender, validation.In(common.UserGenderMale, common.UserGenderFemale)),
		validation.Field(&d.SBPLevel, validation.Min(100.0), validation.Max(179.0)),
		validation.Field(&d.TotalCholesterolLevel, validation.Min(3.0), validation.Max(6.9)),
	)
	if err != nil {
		var errBytes []byte
		errBytes, err = json.Marshal(err)
		if err != nil {
			return err
		}

		var validationErrors map[string]string
		if err = json.Unmarshal(errBytes, &validationErrors); err != nil {
			return err
		}

		if validationError, found := validationErrors["gender"]; found {
			return fmt.Errorf("%w: %v", errors.ErrInvalidGender, validationError)
		}
		if validationError, found := validationErrors["sbpLevel"]; found {
			return fmt.Errorf("%w: %v", errors.ErrInvalidSBPLevel, validationError)
		}
		if validationError, found := validationErrors["totalCholesterolLevel"]; found {
			return fmt.Errorf("%w: %v", errors.ErrInvalidTotalCholesterolLevel, validationError)
		}
		if validationError, found := validationErrors["age"]; found {
			return fmt.Errorf("%w: %v", errors.ErrInvalidAge, validationError)
		}

		return errors.ErrInvalidScoreData
	}
	return nil
}
