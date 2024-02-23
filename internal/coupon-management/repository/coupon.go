package repository

import (
	"context"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/database"
)

// CouponRepository defines the interface for coupon-related database operations.
type CouponRepository interface {
	// Create creates a new coupon in the database and returns its ID.
	Create(ctx context.Context, db database.Executor, data *entity.Coupon) (int64, error)

	// DeleteByID deletes a coupon with the specified ID from the database.
	DeleteByID(ctx context.Context, db database.Executor, id int64) error

	// RetrieveByCode retrieves a coupon from the database based on its unique code.
	RetrieveByCode(ctx context.Context, db database.Executor, code string) (*entity.Coupon, error)
}
