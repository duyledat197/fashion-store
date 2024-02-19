package entity

import "github.com/jackc/pgx/v5/pgtype"

// UsedCouponType ...
type UsedCouponType string

const (
	// UsedCouponTypeUser ...
	UsedCouponTypeUser UsedCouponType = "USER"
	// UsedCouponTypeProduct ...
	UsedCouponTypeProduct UsedCouponType = "PRODUCT"
)

// UsedCoupon ...
type UsedCoupon struct {
	CouponID  pgtype.Int8
	UserID    pgtype.Int8
	Type      UsedCouponType
	CreatedBy pgtype.Int8
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}
