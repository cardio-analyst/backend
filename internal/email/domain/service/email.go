package service

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/cardio-analyst/backend/internal/email/ports/client"
	"github.com/cardio-analyst/backend/pkg/model"
)

type EmailService struct {
	sender client.SMTP
}

func NewEmailService(sender client.SMTP) *EmailService {
	return &EmailService{
		sender: sender,
	}
}

func (s *EmailService) EmailMessagesHandler() func(data []byte) error {
	return func(data []byte) error {
		var message model.SendEmailMessage
		if err := json.Unmarshal(data, &message); err != nil {
			log.Errorf("unmarshalling send email message body: %v", err)
			return err
		}

		if err := s.sender.SendFilePDF(message.Receivers, message.Subject, message.Body, message.FileData); err != nil {
			log.Errorf("sending PDF file: %v", err)
			return err
		}

		return nil
	}
}
