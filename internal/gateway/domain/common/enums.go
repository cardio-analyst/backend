package common

// possible models.BasicIndicators Gender values
const (
	UserGenderMale    = "Мужской"
	UserGenderFemale  = "Женский"
	UserGenderUnknown = "Не выбрано"
)

// possible models.Lifestyle EventsParticipation values
const (
	EventsParticipationFrequently    = "Более 1 раза в неделю"
	EventsParticipationNotFrequently = "1 раз в неделю"
)

// possible models.Lifestyle PhysicalActivity values
const (
	PhysicalActivityOneInWeek         = "Тренировка 1 раз в неделю"
	PhysicalActivityMoreThanOneInWeek = "Тренировка более 1 раза в неделю"
	PhysicalActivityOneInDay          = "Тренировка раз в день"
)

// possible SCORE scale values
const (
	ScaleUnknown  = "unknown"
	ScalePositive = "positive"
	ScaleNeutral  = "neutral"
	ScaleNegative = "negative"
)
