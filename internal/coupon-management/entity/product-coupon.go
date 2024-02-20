package entity

import (
	"database/sql"
)

// ProductCoupon ...
type ProductCoupon struct {
	CouponID  sql.NullInt64 `db:"coupon_id"`
	ProductID sql.NullInt64 `db:"product_id"`
	Used      sql.NullInt64 `db:"used"`
	Total     sql.NullInt64 `db:"total"`
	CreatedBy sql.NullInt64 `db:"created_by"`
	CreatedAt sql.NullTime  `db:"created_at"`
	UpdatedAt sql.NullTime  `db:"updated_at"`
}

func (t *ProductCoupon) TableName() string {
	return "product_coupons"
}
