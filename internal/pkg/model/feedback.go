package model

type Feedback struct {
	ID             uint64   `json:"-" db:"id"`
	UserID         uint64   `json:"userId" db:"user_id"`
	UserFirstName  string   `json:"userFirstName" db:"user_first_name"`
	UserLastName   string   `json:"userLastName" db:"user_last_name"`
	UserMiddleName string   `json:"userMiddleName,omitempty" db:"user_middle_name,omitempty"`
	UserLogin      string   `json:"userLogin" db:"user_login"`
	UserEmail      string   `json:"userEmail" db:"user_email"`
	Mark           int16    `json:"mark" db:"mark"`
	Message        string   `json:"message,omitempty" db:"message,omitempty"`
	CreatedAt      Datetime `json:"createdAt" db:"created_at"`
}
