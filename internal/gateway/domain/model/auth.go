package model

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Session struct {
	UserID       uint64   `json:"-" db:"user_id"`
	RefreshToken string   `json:"-" db:"refresh_token"`
	Whitelist    []string `json:"-" db:"whitelist"`
}

func (s Session) IsIPAllowed(userIP string) bool {
	for _, ip := range s.Whitelist {
		if userIP == ip {
			return true
		}
	}
	return false
}
