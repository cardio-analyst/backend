package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"

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

func (s *FeedbackService) FindAll(criteria model.FeedbackCriteria) ([]model.Feedback, int64, error) {
	feedbacks, err := s.repository.FindAll(criteria)
	if err != nil {
		return nil, 0, err
	}

	feedbacksNum, err := s.repository.Count(criteria)
	if err != nil {
		return nil, 0, err
	}

	var totalPages int64
	if criteria.Limit > 0 {
		limitFloat := float64(criteria.Limit)
		feedbacksNumFloat := float64(feedbacksNum)
		totalPages = int64(math.Ceil(feedbacksNumFloat / limitFloat))
	}

	return feedbacks, totalPages, nil
}

func (s *FeedbackService) ToggleFeedbackViewed(id uint64) error {
	feedback, err := s.repository.One(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrFeedbackNotFound
		}
		return err
	}

	return s.repository.UpdateViewed(feedback.ID, !feedback.Viewed)
}
