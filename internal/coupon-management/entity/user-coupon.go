package entity

import (
	"database/sql"
)

// UserCoupon represents an entity for tracking user-specific coupon information in the database.
type UserCoupon struct {
	CouponID  sql.NullInt64 `db:"coupon_id"`  // ID of the associated coupon
	UserID    sql.NullInt64 `db:"user_id"`    // ID of the user associated with the coupon
	Used      sql.NullInt64 `db:"used"`       // Number of times the coupon has been used
	Total     sql.NullInt64 `db:"total"`      // Total number of coupons available to the user
	CreatedBy sql.NullInt64 `db:"created_by"` // User ID who created the user coupon entry
	CreatedAt sql.NullTime  `db:"created_at"` // User coupon creation timestamp
	UpdatedAt sql.NullTime  `db:"updated_at"` // User coupon last update timestamp
}

// TableName returns the table name for the UserCoupon entity.
func (t *UserCoupon) TableName() string {
	return "user_coupons"
}
