package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

type UsedCouponRepository interface {
	ListUsedCouponByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entity.CouponUsedCoupon, error)
}
