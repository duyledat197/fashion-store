package entity

import (
	"database/sql"
)

// ProductCoupon represents the relationship between products and coupons in the database.
type ProductCoupon struct {
	CouponID  sql.NullInt64 `db:"coupon_id"`  // ID of the associated coupon
	ProductID sql.NullInt64 `db:"product_id"` // ID of the associated product
	Used      sql.NullInt64 `db:"used"`       // Number of times the coupon associated with the product has been used
	Total     sql.NullInt64 `db:"total"`      // Total available coupons for the product
	CreatedBy sql.NullInt64 `db:"created_by"` // User ID who created the product coupon relationship
	CreatedAt sql.NullTime  `db:"created_at"` // Relationship creation timestamp
	UpdatedAt sql.NullTime  `db:"updated_at"` // Relationship last update timestamp
}

// TableName returns the table name for the ProductCoupon entity.
func (t *ProductCoupon) TableName() string {
	return "product_coupons"
}
