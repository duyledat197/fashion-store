package entity

import "database/sql"

// PurchasedProduct represents the structure of a purchased product entity in the database.
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

// TableName returns the name of the database table associated with the PurchasedProduct entity.
func (u *PurchasedProduct) TableName() string {
	return "purchased_products"
}
