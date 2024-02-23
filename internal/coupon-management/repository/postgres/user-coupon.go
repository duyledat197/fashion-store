package postgres

import (
	"context"
	"fmt"
	"strings"

	"trintech/review/internal/coupon-management/entity"
	"trintech/review/internal/coupon-management/repository"
	"trintech/review/pkg/database"
)

// userCouponRepository is an implementation of the UserCouponRepository interface for PostgreSQL.
type userCouponRepository struct{}

// NewUserCouponRepository creates a new instance of the userCouponRepository.
func NewUserCouponRepository() repository.UserCouponRepository {
	return &userCouponRepository{}
}

// Create inserts a new user coupon record into the database.
func (r *userCouponRepository) Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error {
	fieldNames, values := database.FieldMap(data)
	placeHolders := database.GetPlaceholders(len(fieldNames))

	stmt := fmt.Sprintf(`
		INSERT INTO %s("%s")
		VALUES(%s)
	`, data.TableName(), strings.Join(fieldNames, "\",\""), placeHolders)
	if _, err := db.ExecContext(ctx, stmt, values...); err != nil {
		return err
	}

	return nil
}

// DeleteByCouponID deletes user coupons based on the given couponID (not implemented).
func (r *userCouponRepository) DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error {
	panic("not implemented") // TODO: Implement
}

// RetrieveByCouponIDUserID retrieves a user coupon record from the database based on couponID and userID.
func (r *userCouponRepository) RetrieveByCouponIDUserID(ctx context.Context, db database.Executor, couponID, userID int64) (*entity.UserCoupon, error) {
	e := &entity.UserCoupon{}
	fieldNames, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT "%s"
		FROM %s
		WHERE coupon_id = $1 AND user_id = $2
	`, strings.Join(fieldNames, "\",\""), e.TableName())

	if err := db.QueryRowContext(ctx, stmt, &couponID, &userID).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}
