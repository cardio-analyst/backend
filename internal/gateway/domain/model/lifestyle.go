package model

type Lifestyle struct {
	UserID                 uint64 `json:"-" db:"user_id"`
	FamilyStatus           string `json:"familyStatus" db:"family_status"`
	EventsParticipation    string `json:"eventsParticipation" db:"events_participation"`
	PhysicalActivity       string `json:"physicalActivity" db:"physical_activity"`
	WorkStatus             string `json:"workStatus" db:"work_status"`
	SignificantValueHigh   string `json:"significantValueHigh" db:"significant_value_high"`
	SignificantValueMedium string `json:"significantValueMedium" db:"significant_value_medium"`
	SignificantValueLow    string `json:"significantValueLow" db:"significant_value_low"`
}
