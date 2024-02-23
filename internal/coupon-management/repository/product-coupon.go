package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

// ProductCouponRepository defines the interface for product-coupon related database operations.
type ProductCouponRepository interface {
	// Create creates a new product-coupon association in the database.
	Create(ctx context.Context, db database.Executor, data *entity.ProductCoupon) error

	// DeleteByCouponID deletes product-coupon associations by the specified coupon ID.
	DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error

	// RetrieveByCouponID retrieves product-coupon associations by the specified coupon ID.
	RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.ProductCoupon, error)
}
