package service

import (
	"bytes"
	"fmt"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/smtp"
	"html/template"
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

var _ service.EmailService = (*emailService)(nil)

type emailService struct {
	sender smtp.Client
}

func NewEmailService(sender smtp.Client) *emailService {
	return &emailService{
		sender: sender,
	}
}

func (s *emailService) SendReport(receivers []string, reportPath string, userData models.User) error {
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
