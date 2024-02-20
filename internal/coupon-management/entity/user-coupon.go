package entity

import (
	"database/sql"
)

// UserCoupon ...
type UserCoupon struct {
	CouponID  sql.NullInt64 `db:"coupon_id"`
	UserID    sql.NullInt64 `db:"user_id"`
	Used      sql.NullInt64 `db:"used"`
	Total     sql.NullInt64 `db:"total"`
	CreatedBy sql.NullInt64 `db:"created_by"`
	CreatedAt sql.NullTime  `db:"created_at"`
	UpdatedAt sql.NullTime  `db:"updated_at"`
}

func (t *UserCoupon) TableName() string {
	return "user_coupons"
}
