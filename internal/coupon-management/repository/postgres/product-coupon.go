package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/internal/coupon-management/repository"
	"trintech/review/pkg/database"
)

type productCouponRepository struct {
}

func NewProductCouponRepository() repository.ProductCouponRepository {
	return &productCouponRepository{}
}

func (r *productCouponRepository) Create(ctx context.Context, db database.Executor, data *entity.ProductCoupon) error {
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

func (r *productCouponRepository) DeleteByCouponID(ctx context.Context, db database.Executor, couponID int64) error {
	e := &entity.ProductCoupon{}
	stmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE coupon_id = $1
	`, e.TableName())

	result, err := db.ExecContext(ctx, stmt, &couponID)
	if err != nil {
		return err
	}
	rowEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowEffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *productCouponRepository) RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.ProductCoupon, error) {
	e := &entity.ProductCoupon{}
	fieldNames, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = $1
	`, strings.Join(fieldNames, ","), e.TableName())

	if err := db.QueryRowContext(ctx, stmt, &couponID).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}
