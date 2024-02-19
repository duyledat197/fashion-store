package entity

import "github.com/jackc/pgx/v5/pgtype"

// ProductCoupon ...
type ProductCoupon struct {
	CouponID  pgtype.Int8
	ProductID pgtype.Int8
	CreatedBy pgtype.Int8
	Amount    pgtype.Int8
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}
