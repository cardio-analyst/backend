package models

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Diseases struct {
	ID                   uint64 `json:"-" db:"id"`
	UserID               uint64 `json:"-" db:"user_id"`
	CvdsPredisposition   string `json:"cvdsPredisposition" db:"cvds_predisposition"`
	TakeStatins          bool   `json:"takeStatins" db:"take_statins"`
	Ckd                  bool   `json:"ckd" db:"ckd"`
	ArterialHypertension bool   `json:"arterial_hypertension" db:"arterial_hypertension"`
	CardiacIschemia      bool   `json:"cardiacIschemia" db:"cardiac_ischemia"`
	TypeTwoDiabets       bool   `json:"typeTwoDiabets" db:"type_two_diabets"`
	InfarctionOrStroke   string `json:"infarctionOrStroke" db:"infarction_or_stroke"`
	Atherosclerosis      bool   `json:"atherosclerosis" db:"atherosclerosis"`
	OtherCvdsDiseases    string `json:"otherCvdsDiseases" db:"other_cvds_diseases"`
}

func (d Diseases) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CvdsPredisposition, validation.Required),
		validation.Field(&d.TakeStatins, validation.Required),
		validation.Field(&d.ArterialHypertension, validation.Required),
		validation.Field(&d.CardiacIschemia, validation.Required),
		validation.Field(&d.TypeTwoDiabets, validation.Required),
		validation.Field(&d.InfarctionOrStroke, validation.Required),
		validation.Field(&d.Atherosclerosis, validation.Required),
	)
}
