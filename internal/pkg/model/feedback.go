package model

type Feedback struct {
	ID             uint64   `json:"-" db:"id"`
	UserID         uint64   `json:"user_id" db:"user_id"`
	UserFirstName  string   `json:"user_first_name" db:"user_first_name"`
	UserLastName   string   `json:"user_last_name" db:"user_last_name"`
	UserMiddleName string   `json:"user_middle_name,omitempty" db:"user_middle_name,omitempty"`
	UserLogin      string   `json:"user_login" db:"user_login"`
	UserEmail      string   `json:"user_email" db:"user_email"`
	Mark           int16    `json:"mark" db:"mark"`
	Message        string   `json:"message,omitempty" db:"message,omitempty"`
	CreatedAt      Datetime `json:"createdAt" db:"created_at"`
}
