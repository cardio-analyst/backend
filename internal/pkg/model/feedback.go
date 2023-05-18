package model

import "errors"

var ErrFeedbackNotFound = errors.New("feedback not found")

type Feedback struct {
	ID             uint64   `json:"id" db:"id"`
	UserID         uint64   `json:"userId" db:"user_id"`
	UserFirstName  string   `json:"userFirstName" db:"user_first_name"`
	UserLastName   string   `json:"userLastName" db:"user_last_name"`
	UserMiddleName string   `json:"userMiddleName,omitempty" db:"user_middle_name,omitempty"`
	UserLogin      string   `json:"userLogin" db:"user_login"`
	UserEmail      string   `json:"userEmail" db:"user_email"`
	Mark           int16    `json:"mark" db:"mark"`
	Message        string   `json:"message,omitempty" db:"message,omitempty"`
	Version        string   `json:"version" db:"version"`
	Viewed         bool     `json:"viewed" db:"viewed"`
	CreatedAt      Datetime `json:"createdAt" db:"created_at"`
}

type OrderingType int64

const (
	OrderingTypeDisabled OrderingType = iota
	OrderingTypeASC
	OrderingTypeDESC
)

type FeedbackCriteria struct {
	MarkOrdering    OrderingType
	VersionOrdering OrderingType
	Viewed          *bool
	Limit           int64
	Page            int64
}
