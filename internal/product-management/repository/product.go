package repository

import (
	"context"

	"trintech/review/internal/product-management/entity"
	"trintech/review/pkg/database"
)

type ProductRepository interface {
	List(ctx context.Context, db database.Executor, offset, limit int64) ([]*entity.Product, error)
	RetrieveByID(ctx context.Context, db database.Executor, id int64) (*entity.Product, error)
	Create(ctx context.Context, db database.Executor, data *entity.Product) (int64, error)
	UpdateByID(ctx context.Context, db database.Executor, id int64, data *entity.Product) error
	DeleteByID(ctx context.Context, db database.Executor, id int64) error
	DeleteByIDs(ctx context.Context, db database.Executor, ids []int64) error
	Count(ctx context.Context, db database.Executor) (int64, error)
}
