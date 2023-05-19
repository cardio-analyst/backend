package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const (
	reportSubject  = "Отчёт по показателям здоровья пациента"
	reportHTMLBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>{{ .title }}</title>
	</head>
	<body>
        <p>{{ .firstName }} {{ .lastName }}, {{ .birthDate }} г.р.</p>
		<p>Сгенерирован сервисом "Кардио Аналитик".</p>
	</body>
</html>`
)

var _ service.EmailService = (*EmailService)(nil)

type EmailService struct {
	emailsPublisher client.Publisher
}

func NewEmailService(emailsPublisher client.Publisher) *EmailService {
	return &EmailService{
		emailsPublisher: emailsPublisher,
	}
}

func (s *EmailService) SendReport(receivers []string, reportFilePath string, userData model.User) error {
	reportTemplate, err := template.New("report").Parse(reportHTMLBody)
	if err != nil {
		return err
	}

	reportBodyBuffer := &bytes.Buffer{}
	if err = reportTemplate.Execute(reportBodyBuffer, map[string]interface{}{
		"title":     fmt.Sprintf("%v %v %v", reportSubject, userData.FirstName, userData.LastName),
		"firstName": userData.FirstName,
		"lastName":  userData.LastName,
		"birthDate": userData.BirthDate.String(),
	}); err != nil {
		return err
	}

	file, err := os.ReadFile(reportFilePath)
	if err != nil {
		return fmt.Errorf("reading report file: %w", err)
	}

	message := &model.MessageReportEmail{
		Subject:   reportSubject,
		Receivers: receivers,
		Body:      reportBodyBuffer.String(),
		FileData:  file,
	}

	rmqMessageRaw, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("serializing RMQ message: %v", err)
	}

	return s.emailsPublisher.Publish(rmqMessageRaw)
}
