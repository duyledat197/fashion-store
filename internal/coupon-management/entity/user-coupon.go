package entity

import "github.com/jackc/pgx/v5/pgtype"

// UserCoupon ...
type UserCoupon struct {
	CouponID  pgtype.Int8
	UserID    pgtype.Int8
	Amount    pgtype.Int8
	CreatedBy pgtype.Int8
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}
