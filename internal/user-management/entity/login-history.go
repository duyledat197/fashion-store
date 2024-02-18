package entity

import "time"

type LoginHistory struct {
	UserID      int64     `json:"user_id,omitempty"`
	IP          string    `json:"ip,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	LoginAt     time.Time `json:"login_at,omitempty"`
	LogoutAt    time.Time `json:"log_out_at,omitempty"`
}

func (u *LoginHistory) TableName() string {
	return "login_histories"
}
