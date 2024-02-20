package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

type ProductCouponRepository interface {
	Create(ctx context.Context, db database.Executor, data *entity.ProductCoupon) error
	DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
	RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.ProductCoupon, error)
}
