package postgres

import (
	"context"
	"fmt"
	"strings"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/internal/coupon-management/repository"
	"trintech/review/pkg/database"
)

type usedCouponRepository struct {
}

func NewUsedCouponRepository() repository.UsedCouponRepository {
	return &usedCouponRepository{}
}

func (r *usedCouponRepository) ListUsedCouponByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entity.CouponUsedCoupon, error) {
	e := &entity.UsedCoupon{}
	cE := &entity.Coupon{}
	fieldNames, _ := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT uc.%s, c.%s
		FROM %s uc
		JOIN %s c
		ON uc.coupon_id = c.id
		WHERE user_id = $1

	`,
		strings.Join(fieldNames, ",uc."),
		strings.Join(fieldNames, ",c."),
		e.TableName(),
		cE.TableName(),
	)

	rows, err := db.QueryContext(ctx, stmt, &userID)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var result []*entity.CouponUsedCoupon
	for rows.Next() {
		var (
			uVal entity.UsedCoupon
			cVal entity.Coupon
		)

		_, uValues := database.FieldMap(&uVal)
		_, cValues := database.FieldMap(&cVal)
		var values []any
		values = append(values, uValues...)
		values = append(values, cValues...)

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		result = append(result, &entity.CouponUsedCoupon{
			UsedCoupon: &uVal,
			Coupon:     &cVal,
		})
	}

	return result, nil
}

func (r *usedCouponRepository) Create(ctx context.Context, db database.Executor, data *entity.UsedCoupon) error {
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
