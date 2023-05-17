package service

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/cardio-analyst/backend/internal/analytics/ports/client"
	"github.com/cardio-analyst/backend/internal/analytics/ports/storage"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type FeedbackService struct {
	repository storage.FeedbackRepository
	consumer   client.FeedbackConsumer
}

func NewFeedbackService(repository storage.FeedbackRepository, consumer client.FeedbackConsumer) *FeedbackService {
	return &FeedbackService{
		repository: repository,
		consumer:   consumer,
	}
}

func (s *FeedbackService) ListenToFeedbackMessages() error {
	handler := s.feedbackMessagesHandler()

	if err := s.consumer.Consume(handler); err != nil {
		return fmt.Errorf("consuming feedback messages: %w", err)
	}

	return nil
}

func (s *FeedbackService) feedbackMessagesHandler() func(data []byte) error {
	return func(data []byte) error {
		var message model.MessageFeedback
		if err := json.Unmarshal(data, &message); err != nil {
			log.Errorf("unmarshalling feedback message body: %v", err)
			return err
		}

		feedback := model.Feedback{
			UserID:         message.UserID,
			UserFirstName:  message.UserFirstName,
			UserLastName:   message.UserLastName,
			UserMiddleName: message.UserMiddleName,
			UserLogin:      message.UserLogin,
			UserEmail:      message.UserEmail,
			Mark:           message.Mark,
			Message:        message.Message,
			Version:        message.Version,
		}

		if err := s.repository.Create(feedback); err != nil {
			log.Errorf("creating feedback record: %v", err)
			return err
		}

		return nil
	}
}

func (s *FeedbackService) FindAll() ([]model.Feedback, error) {
	return s.repository.FindAll()
}
