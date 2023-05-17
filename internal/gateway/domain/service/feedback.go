package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type FeedbackService struct {
	publisher       client.FeedbackPublisher
	analyticsClient client.Analytics
}

func NewFeedbackService(publisher client.FeedbackPublisher, analyticsClient client.Analytics) *FeedbackService {
	return &FeedbackService{
		publisher:       publisher,
		analyticsClient: analyticsClient,
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

	return s.publisher.Publish(rmqMessageRaw)
}

func (s *FeedbackService) FindAll() ([]model.Feedback, error) {
	return s.analyticsClient.FindAllFeedbacks(context.TODO())
}
