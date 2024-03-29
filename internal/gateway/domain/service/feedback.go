package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type FeedbackService struct {
	feedbackPublisher client.Publisher
	analyticsClient   client.Analytics
}

func NewFeedbackService(feedbackPublisher client.Publisher, analyticsClient client.Analytics) *FeedbackService {
	return &FeedbackService{
		feedbackPublisher: feedbackPublisher,
		analyticsClient:   analyticsClient,
	}
}

func (s *FeedbackService) Send(mark int16, text, version string, user model.User) error {
	message := &model.MessageFeedback{
		UserID:         user.ID,
		UserFirstName:  user.FirstName,
		UserLastName:   user.LastName,
		UserMiddleName: user.MiddleName,
		UserLogin:      user.Login,
		UserEmail:      user.Email,
		Mark:           mark,
		Message:        text,
		Version:        version,
	}

	rmqMessageRaw, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("serializing RMQ message: %v", err)
	}

	return s.feedbackPublisher.Publish(rmqMessageRaw)
}

func (s *FeedbackService) FindAll(criteria model.FeedbackCriteria) ([]model.Feedback, int64, error) {
	return s.analyticsClient.FindAllFeedbacks(context.TODO(), criteria)
}

func (s *FeedbackService) ToggleFeedbackViewed(id uint64) error {
	return s.analyticsClient.ToggleFeedbackViewed(context.TODO(), id)
}
