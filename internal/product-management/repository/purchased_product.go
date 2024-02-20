package repository

import (
	"context"

	"trintech/review/internal/product-management/entity"
	"trintech/review/pkg/database"
)

type PurchasedProductRepository interface {
	Create(ctx context.Context, db database.Executor, data *entity.PurchasedProduct) error
}
