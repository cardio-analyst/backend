package storage

// Storage represents the data storage layer (database).
type Storage interface {
	Feedback() FeedbackRepository
	RegionUsers() RegionUsersRepository
}
