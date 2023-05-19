package service

type StatisticsService interface {
	ListenToRegistrationMessages() error
	AllUsersByRegions() (map[string]int64, error)
}
