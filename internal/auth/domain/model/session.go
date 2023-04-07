package model

import "errors"

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID       uint64   `bson:"user_id,omitempty"`
	RefreshToken string   `bson:"refresh_token"`
	Whitelist    []string `bson:"whitelist"`
}

func (s Session) IsIPAllowed(userIP string) bool {
	for _, ip := range s.Whitelist {
		if userIP == ip {
			return true
		}
	}
	return false
}
