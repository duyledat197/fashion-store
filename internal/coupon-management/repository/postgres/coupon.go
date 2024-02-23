package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"trintech/review/internal/coupon-management/entity"
	repository "trintech/review/internal/coupon-management/repository"
	"trintech/review/pkg/database"
)

// couponRepository is an implementation of the CouponRepository interface for PostgreSQL.
type couponRepository struct{}

// NewCouponRepository creates a new instance of the couponRepository.
func NewCouponRepository() repository.CouponRepository {
	return &couponRepository{}
}

// Create inserts a new coupon record into the database and returns the generated ID.
func (r *couponRepository) Create(ctx context.Context, db database.Executor, data *entity.Coupon) (int64, error) {
	// Get field names and values excluding the "id" field.
	fieldNames, values := database.FieldMap(data)
	fieldNames = fieldNames[1:]
	values = values[1:]
	placeHolders := database.GetPlaceholders(len(fieldNames))

	stmt := fmt.Sprintf(`
		INSERT INTO %s("%s")
		VALUES(%s)
		RETURNING id
	`, data.TableName(), strings.Join(fieldNames, "\",\""), placeHolders)
	var id int64
	log.Println(stmt)
	if err := db.QueryRowContext(ctx, stmt, values...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteByID deletes a coupon record from the database based on its ID.
func (r *couponRepository) DeleteByID(ctx context.Context, db database.Executor, id int64) error {
	e := &entity.Coupon{}
	stmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, e.TableName())

	result, err := db.ExecContext(ctx, stmt, &id)
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

// RetrieveByCode retrieves a coupon record from the database based on its code.
func (r *couponRepository) RetrieveByCode(ctx context.Context, db database.Executor, code string) (*entity.Coupon, error) {
	e := &entity.Coupon{}
	fieldNames, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT "%s"
		FROM %s
		WHERE code = $1
	`, strings.Join(fieldNames, "\",\""), e.TableName())

	if err := db.QueryRowContext(ctx, stmt, &code).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}
