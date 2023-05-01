package service

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/cardio-analyst/backend/internal/email/ports/client"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

type EmailService struct {
	sender   client.SMTP
	consumer client.EmailConsumer
}

func NewEmailService(sender client.SMTP, consumer client.EmailConsumer) *EmailService {
	return &EmailService{
		sender:   sender,
		consumer: consumer,
	}
}

func (s *EmailService) ListenToEmailMessages() error {
	handler := s.emailMessagesHandler()

	if err := s.consumer.Consume(handler); err != nil {
		return fmt.Errorf("consuming email messages: %w", err)
	}

	return nil
}

func (s *EmailService) emailMessagesHandler() func(data []byte) error {
	return func(data []byte) error {
		var message model.MessageReportEmail
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
