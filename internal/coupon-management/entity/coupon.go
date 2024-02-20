// Package entity ...
package entity

import (
	"database/sql"
)

// Coupon ...
type Coupon struct {
	ID           sql.NullInt64   `db:"id"`
	Code         sql.NullString  `db:"code"`
	From         sql.NullTime    `db:"from"`
	To           sql.NullTime    `db:"to"`
	Used         sql.NullInt64   `db:"used"`
	Total        sql.NullInt64   `db:"total"`
	Type         sql.NullString  `db:"coupon_type"`
	Value        sql.NullFloat64 `db:"value"`
	ImageURL     sql.NullString  `db:"image_url"`
	Description  sql.NullString  `db:"description"`
	DiscountType sql.NullString  `db:"discount_type"`
	CreatedBy    sql.NullInt64   `db:"created_by"`
	CreatedAt    sql.NullTime    `db:"created_at"`
	UpdatedAt    sql.NullTime    `db:"updated_at"`
}

func (t *Coupon) TableName() string {
	return "coupons"
}
