package service

type ReportService interface {
	GenerateReport(userID uint64) (reportPath string, err error)
}
