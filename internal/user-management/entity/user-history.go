package entity

import "time"

type UserHistory struct {
	UserID    string    `json:"user_id,omitempty"`
	URL       string    `json:"url,omitempty"`
	Method    string    `json:"method,omitempty"`
	SessionID string    `json:"session_id,omitempty"`
	ActionAt  time.Time `json:"action_at,omitempty"`
}
