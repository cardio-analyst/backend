package models

type Recommendation struct {
	What string `json:"what"`
	Why  string `json:"why"`
	How  string `json:"how"`
}
