package service

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/pkg/model"
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
	sender client.SMTP
}

func NewEmailService(sender client.SMTP) *EmailService {
	return &EmailService{
		sender: sender,
	}
}

func (s *EmailService) SendReport(receivers []string, reportPath string, userData model.User) error {
	reportTemplate, err := template.New("report").Parse(reportHTMLBody)
	if err != nil {
		return err
	}

	reportBodyBuffer := &bytes.Buffer{}
	if err = reportTemplate.Execute(reportBodyBuffer, map[string]interface{}{
		"title":      fmt.Sprintf("%v %v %v", reportSubject, userData.FirstName, userData.LastName),
		"firstName":  userData.FirstName,
		"lastName":   userData.LastName,
		"birthDate":  userData.BirthDate.String(),
		"reportPath": reportPath,
	}); err != nil {
		return err
	}

	return s.sender.SendFile(receivers, reportSubject, reportBodyBuffer.String(), reportPath)
}
