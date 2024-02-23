package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

// UsedCouponRepository defines the interface for used coupon related database operations.
type UsedCouponRepository interface {
	// ListUsedCouponByUserID retrieves a list of used coupons associated with a specific user ID.
	ListUsedCouponByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entity.CouponUsedCoupon, error)

	// Create creates a new entry for a used coupon in the database.
	Create(ctx context.Context, db database.Executor, data *entity.UsedCoupon) error
}
