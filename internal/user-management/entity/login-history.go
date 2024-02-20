package entity

import (
	"database/sql"
)

type LoginHistory struct {
	UserID      sql.NullInt64  `db:"user_id"`
	IP          sql.NullString `db:"ip"`
	UserAgent   sql.NullString `db:"user_agent"`
	AccessToken sql.NullString `db:"access_token"`
	LoginAt     sql.NullTime   `db:"login_at"`
	LogoutAt    sql.NullTime   `db:"logout_at"`
}

func (u *LoginHistory) TableName() string {
	return "login_histories"
}
