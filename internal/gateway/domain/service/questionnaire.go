package service

import (
	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

// check whether QuestionnaireService structure implements the service.QuestionnaireService interface
var _ service.QuestionnaireService = (*QuestionnaireService)(nil)

// QuestionnaireService implements service.QuestionnaireService interface.
type QuestionnaireService struct {
	repository storage.QuestionnaireRepository
}

func NewQuestionnaireService(repository storage.QuestionnaireRepository) *QuestionnaireService {
	return &QuestionnaireService{
		repository: repository,
	}
}

func (s *QuestionnaireService) Get(userID uint64) (*domain.Questionnaire, error) {
	questionnaire, err := s.repository.Get(userID)
	if err != nil {
		return nil, err
	}

	return questionnaire, nil
}

func (s *QuestionnaireService) UpdateAnginaRose(questionnaire domain.Questionnaire) error {
	return s.repository.UpdateAnginaRose(questionnaire)
}

func (s *QuestionnaireService) UpdateTreatmentAdherence(questionnaire domain.Questionnaire) error {
	return s.repository.UpdateTreatmentAdherence(questionnaire)
}
