package entity

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRole string

const (
	UserRole_User       = "USER"
	UserRole_Admin      = "ADMIN"
	UserRole_SuperAdmin = "SUPER_ADMIN"
)

type User struct {
	ID        int64              `json:"id,omitempty"`
	UserName  pgtype.Text        `json:"user_name,omitempty"`
	Email     pgtype.Text        `json:"email,omitempty"`
	Password  pgtype.Text        `json:"password,omitempty"`
	Name      pgtype.Text        `json:"name,omitempty"`
	Role      UserRole           `json:"role,omitempty"`
	CreatedAt pgtype.Timestamptz `json:"created_at,omitempty"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}
