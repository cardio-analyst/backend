package service

type ReportService interface {
	GenerateReport(userID uint64) (reportFilePath string, err error)
}
