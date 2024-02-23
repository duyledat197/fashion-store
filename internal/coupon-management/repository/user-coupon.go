package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

// UserCouponRepository defines the interface for user coupon related database operations.
type UserCouponRepository interface {
	// Create creates a new entry for a user coupon in the database.
	Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error

	// DeleteByCouponID deletes user coupons associated with a specific coupon ID from the database.
	DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error

	// RetrieveByCouponIDUserID retrieves a user coupon based on coupon ID and user ID.
	RetrieveByCouponIDUserID(ctx context.Context, db database.Executor, couponID, userID int64) (*entity.UserCoupon, error)
}
