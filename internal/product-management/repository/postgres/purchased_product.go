package postgres

import (
	"context"
	"fmt"
	"strings"

	"trintech/review/internal/product-management/entity"
	"trintech/review/pkg/database"
)

type purchasedProductRepository struct {
}

func (r *purchasedProductRepository) Create(ctx context.Context, db database.Executor, data *entity.PurchasedProduct) error {
	fieldNames, values := database.FieldMap(data)
	placeHolders := database.GetPlaceholders(len(fieldNames))

	stmt := fmt.Sprintf(`
		INSERT INTO %s(%s)
		VALUES(%s)
		RETURNING id
	`, data.TableName(), strings.Join(fieldNames, ","), placeHolders)
	var id int64

	if err := db.QueryRowContext(ctx, stmt, values...).Scan(&id); err != nil {
		return err
	}

	return nil
}
