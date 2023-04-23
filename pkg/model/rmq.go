package model

type MessageReportEmail struct {
	Subject   string   `json:"subject"`
	Receivers []string `json:"receivers"`
	Body      string   `json:"body"`
	FileData  []byte   `json:"file_data"`
}
