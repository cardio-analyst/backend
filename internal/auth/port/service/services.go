package service

// Services represents a layer of business logic.
type Services interface {
	// Auth contains the methods of business logic of working with authorization and authentication.
	Auth() AuthService
	// User contains the methods of business logic of working with users.
	User() UserService
	// Validation contains the methods of business logic of working with data validation.
	Validation() ValidationService
}
