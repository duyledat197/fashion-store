package entity

import (
	"database/sql"
)

// UsedCouponType represents the type of a used coupon (USER, PRODUCT, LIMITED).
type UsedCouponType string

const (
	// UsedCouponTypeUser represents a used coupon associated with a user.
	UsedCouponTypeUser UsedCouponType = "USER"
	// UsedCouponTypeProduct represents a used coupon associated with a product.
	UsedCouponTypeProduct UsedCouponType = "PRODUCT"
	// UsedCouponTypeLimited represents a used coupon with a limited type.
	UsedCouponTypeLimited UsedCouponType = "LIMITED"
)

// UsedCoupon represents an entity for tracking used coupons in the database.
type UsedCoupon struct {
	CouponID  sql.NullInt64 `db:"coupon_id"`  // ID of the associated coupon
	UserID    sql.NullInt64 `db:"user_id"`    // ID of the user associated with the used coupon
	CreatedBy sql.NullInt64 `db:"created_by"` // User ID who created the used coupon entry
	CreatedAt sql.NullTime  `db:"created_at"` // Used coupon creation timestamp
	UpdatedAt sql.NullTime  `db:"updated_at"` // Used coupon last update timestamp
}

// TableName returns the table name for the UsedCoupon entity.
func (t *UsedCoupon) TableName() string {
	return "used_coupons"
}

// CouponUsedCoupon is a struct that represents the combination of a UsedCoupon and Coupon.
// It's used to combine information about a used coupon and the corresponding coupon.
type CouponUsedCoupon struct {
	UsedCoupon *UsedCoupon // Information about the used coupon
	Coupon     *Coupon     // Information about the associated coupon
}
