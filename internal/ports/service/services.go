package service

// Services represents a layer of business logic.
type Services interface {
	// User contains the methods of business logic of working with users.
	User() UserService
	// Auth contains the methods of business logic of working with authorization.
	Auth() AuthService
}
