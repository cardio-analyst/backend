package service

import (
	"fmt"
	"os"

	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/smtp"
)

var _ service.ReportService = (*pdfReportService)(nil)

type pdfReportService struct {
	sender smtp.Client
}

func NewPDFReportService(sender smtp.Client) *pdfReportService {
	return &pdfReportService{
		sender: sender,
	}
}

func (s *pdfReportService) GenerateReport(userID uint64) (string, error) {
	if err := os.WriteFile("report.txt", []byte(fmt.Sprintf("Hello, %v", userID)), 0755); err != nil {
		return "", err
	}
	return "report.txt", nil
}
