// Package entity ...
package entity

import (
	"database/sql"
)

// Coupon represents the database entity for coupons.
type Coupon struct {
	ID           sql.NullInt64   `db:"id"`            // Coupon ID
	Code         sql.NullString  `db:"code"`          // Coupon code
	From         sql.NullTime    `db:"from"`          // Validity start time
	To           sql.NullTime    `db:"to"`            // Validity end time
	Used         sql.NullInt64   `db:"used"`          // Number of times the coupon has been used
	Total        sql.NullInt64   `db:"total"`         // Total available coupons
	Type         sql.NullString  `db:"coupon_type"`   // Type of coupon
	Value        sql.NullFloat64 `db:"value"`         // Value or percentage of the coupon discount
	ImageURL     sql.NullString  `db:"image_url"`     // URL of the coupon image
	Description  sql.NullString  `db:"description"`   // Coupon description
	DiscountType sql.NullString  `db:"discount_type"` // Type of discount (e.g., fixed value or percentage)
	CreatedBy    sql.NullInt64   `db:"created_by"`    // User ID who created the coupon
	CreatedAt    sql.NullTime    `db:"created_at"`    // Coupon creation timestamp
	UpdatedAt    sql.NullTime    `db:"updated_at"`    // Coupon last update timestamp
}

// TableName returns the table name for the Coupon entity.
func (t *Coupon) TableName() string {
	return "coupons"
}
