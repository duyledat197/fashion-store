package entity

import (
	"database/sql"
)

type UserRole string

const (
	UserRole_User       = "USER"
	UserRole_Admin      = "ADMIN"
	UserRole_SuperAdmin = "SUPER_ADMIN"
)

type User struct {
	ID        sql.NullInt64  `db:"id"`
	UserName  sql.NullString `db:"user_name"`
	Email     sql.NullString `db:"email"`
	Password  sql.NullString `db:"password"`
	Name      sql.NullString `db:"name"`
	Role      UserRole       `db:"role"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
