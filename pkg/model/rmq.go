package model

type SendEmailMessage struct {
	Subject   string   `json:"subject"`
	Receivers []string `json:"receivers"`
	Body      string   `json:"body"`
	FileData  []byte   `json:"file_data"`
}
