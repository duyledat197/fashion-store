package postgres

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/internal/coupon-management/repository"
	"trintech/review/pkg/database"
)

type userCouponRepository struct {
}

func NewUserCouponRepository() repository.UserCouponRepository {
	return &userCouponRepository{}
}

func (r *userCouponRepository) Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error {
	panic("not implemented") // TODO: Implement
}

func (r *userCouponRepository) DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error {
	panic("not implemented") // TODO: Implement
}

func (r *userCouponRepository) RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.UserCoupon, error) {
	panic("not implemented") // TODO: Implement
}
