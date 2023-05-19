package model

type MessageRegistration struct {
	Region string `json:"user_region"`
}

type MessageReportEmail struct {
	Subject   string   `json:"subject"`
	Receivers []string `json:"receivers"`
	Body      string   `json:"body"`
	FileData  []byte   `json:"file_data"`
}

type MessageFeedback struct {
	UserID         uint64 `json:"user_id"`
	UserFirstName  string `json:"user_first_name"`
	UserLastName   string `json:"user_last_name"`
	UserMiddleName string `json:"user_middle_name,omitempty"`
	UserLogin      string `json:"user_login"`
	UserEmail      string `json:"user_email"`
	Mark           int16  `json:"mark"`
	Message        string `json:"message,omitempty"`
	Version        string `json:"version"`
}
