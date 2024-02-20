package postgres

import (
	"context"
	"fmt"
	"strings"

	"trintech/review/internal/product-management/entity"
	"trintech/review/internal/product-management/repository"
	"trintech/review/pkg/database"
)

type purchasedProductRepository struct {
}

func NewPurchasedProductRepository() repository.PurchasedProductRepository {
	return &purchasedProductRepository{}
}

func (r *purchasedProductRepository) Create(ctx context.Context, db database.Executor, data *entity.PurchasedProduct) error {
	fieldNames, values := database.FieldMap(data)
	placeHolders := database.GetPlaceholders(len(fieldNames))

	stmt := fmt.Sprintf(`
		INSERT INTO %s(%s)
		VALUES(%s)
	`, data.TableName(), strings.Join(fieldNames, ","), placeHolders)

	if _, err := db.ExecContext(ctx, stmt, values...); err != nil {
		return err
	}

	return nil
}
