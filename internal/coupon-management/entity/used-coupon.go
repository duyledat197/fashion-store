package entity

import (
	"database/sql"
)

// UsedCouponType ...
type UsedCouponType string

const (
	// UsedCouponTypeUser ...
	UsedCouponTypeUser UsedCouponType = "USER"
	// UsedCouponTypeProduct ...
	UsedCouponTypeProduct UsedCouponType = "PRODUCT"
	// UsedCouponTypeLimited ...
	UsedCouponTypeLimited UsedCouponType = "LIMITED"
)

// UsedCoupon ...
type UsedCoupon struct {
	CouponID  sql.NullInt64  `db:"coupon_id"`
	UserID    sql.NullInt64  `db:"user_id"`
	Type      UsedCouponType `db:"type"`
	CreatedBy sql.NullInt64  `db:"created_by"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

func (t *UsedCoupon) TableName() string {
	return "used_coupons"
}

type CouponUsedCoupon struct {
	UsedCoupon *UsedCoupon
	Coupon     *Coupon
}
