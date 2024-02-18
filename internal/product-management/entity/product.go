package entity

import "github.com/jackc/pgx/v5/pgtype"

type Product struct {
	Name        pgtype.Text `json:"name,omitempty"`
	Type        pgtype.Text `json:"type,omitempty"`
	ImageURL    pgtype.Text `json:"image,omitempty"`
	Description pgtype.Text `json:"description,omitempty"`
}

func (u *Product) TableName() string {
	return "products"
}
