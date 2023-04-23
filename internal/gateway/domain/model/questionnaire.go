package model

type Questionnaire struct {
	UserID                  uint64  `json:"-" db:"user_id"`
	AnginaScore             int8    `json:"anginaScore" db:"angina_score"`
	AdherenceDrugTherapy    float64 `json:"adherenceDrugTherapy" db:"adherence_drug_therapy"`
	AdherenceMedicalSupport float64 `json:"adherenceMedicalSupport" db:"adherence_medical_support"`
	AdherenceLifestyleMod   float64 `json:"adherenceLifestyleMod" db:"adherence_lifestyle_mod"`
}
