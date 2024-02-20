package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

type CouponRepository interface {
	Create(ctx context.Context, db database.Executor, data *entity.Coupon) (int64, error)
	DeleteByID(ctx context.Context, db database.Executor, id int64) error
	RetrieveByCode(ctx context.Context, db database.Executor, code string) (*entity.Coupon, error)
}
