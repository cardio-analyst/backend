package service

// Services represents a layer of business logic.
type Services interface {
	Feedback() FeedbackService
	Statistics() StatisticsService
}
