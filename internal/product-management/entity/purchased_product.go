package entity

import (
	"database/sql"
)

// PurchasedProduct ...
type PurchasedProduct struct {
	ProductID sql.NullInt64   `db:"product_id"`
	UserID    sql.NullInt64   `db:"user_id"`
	Price     sql.NullFloat64 `db:"price"`
	Discount  sql.NullFloat64 `db:"discount"`
	Total     sql.NullFloat64 `db:"purchase_total"`
	Coupon    sql.NullString  `db:"apply_coupon"`
	CreatedAt sql.NullTime    `db:"created_at"`
	UpdatedAt sql.NullTime    `db:"updated_at"`
}

func (u *PurchasedProduct) TableName() string {
	return "purchased_products"
}
