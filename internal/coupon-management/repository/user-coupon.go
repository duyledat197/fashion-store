package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

type UserCouponRepository interface {
	Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error
	DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
	RetrieveByCouponIDUserID(ctx context.Context, db database.Executor, couponID, userID int64) (*entity.UserCoupon, error)
}
