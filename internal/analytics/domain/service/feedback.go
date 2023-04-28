package service

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	domain "github.com/cardio-analyst/backend/internal/analytics/domain/model"
	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
	"github.com/cardio-analyst/backend/pkg/model"
)

type FeedbackService struct {
	repository storage.FeedbackRepository
}

func NewFeedbackService(repository storage.FeedbackRepository) *FeedbackService {
	return &FeedbackService{
		repository: repository,
	}
}

func (s *FeedbackService) MessagesHandler() func(data []byte) error {
	return func(data []byte) error {
		var message model.MessageFeedback
		if err := json.Unmarshal(data, &message); err != nil {
			log.Errorf("unmarshalling feedback message body: %v", err)
			return err
		}

		feedback := domain.Feedback{
			UserID:         message.UserID,
			UserFirstName:  message.UserFirstName,
			UserLastName:   message.UserLastName,
			UserMiddleName: message.UserMiddleName,
			UserLogin:      message.UserLogin,
			UserEmail:      message.UserEmail,
			Mark:           message.Mark,
			Message:        message.Message,
		}

		if err := s.repository.Create(feedback); err != nil {
			log.Errorf("creating feedback record: %v", err)
			return err
		}

		return nil
	}
}
