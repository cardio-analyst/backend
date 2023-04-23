package storage

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

type QuestionnaireRepository interface {
	Get(userID uint64) (questionnaire *domain.Questionnaire, err error)
	UpdateAnginaRose(questionnaire domain.Questionnaire) (err error)
	UpdateTreatmentAdherence(questionnaire domain.Questionnaire) (err error)
}
