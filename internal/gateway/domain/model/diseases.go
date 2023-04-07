package model

import "errors"

var ErrUserDiseasesNotFound = errors.New("user diseases record not found")

type Diseases struct {
	UserID                  uint64 `json:"-" db:"user_id"`
	CVDPredisposed          bool   `json:"cvdPredisposed" db:"cvd_predisposed"`
	TakesStatins            bool   `json:"takesStatins" db:"takes_statins"`
	HasChronicKidneyDisease bool   `json:"hasChronicKidneyDisease" db:"has_chronic_kidney_disease"`
	HasArterialHypertension bool   `json:"hasArterialHypertension" db:"has_arterial_hypertension"`
	HasIschemicHeartDisease bool   `json:"hasIschemicHeartDisease" db:"has_ischemic_heart_disease"`
	HasTypeTwoDiabetes      bool   `json:"hasTypeTwoDiabetes" db:"has_type_two_diabetes"`
	HadInfarctionOrStroke   bool   `json:"hadInfarctionOrStroke" db:"had_infarction_or_stroke"`
	HasAtherosclerosis      bool   `json:"hasAtherosclerosis" db:"has_atherosclerosis"`
	HasOtherCVD             bool   `json:"hasOtherCVD" db:"has_other_cvd"`
}
