package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"trintech/review/internal/product-management/entity"
	"trintech/review/internal/product-management/repository"
	"trintech/review/pkg/database"
)

type productRepository struct{}

func NewProductRepository() repository.ProductRepository {
	return &productRepository{}
}

func (r *productRepository) List(ctx context.Context, db database.Executor, offset int64, limit int64) ([]*entity.Product, error) {
	e := &entity.Product{}
	fieldNames, _ := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT %s
		FROM %s
		LIMIT $1
		OFFSET $2
	`, strings.Join(fieldNames, ","), e.TableName())

	rows, err := db.QueryContext(ctx, stmt, &limit, &offset)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var result []*entity.Product
	for rows.Next() {
		var val entity.Product
		_, values := database.FieldMap(&val)
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		result = append(result, &val)
	}

	return result, nil
}

func (r *productRepository) RetrieveByID(ctx context.Context, db database.Executor, id int64) (*entity.Product, error) {
	e := &entity.Product{}
	fieldNames, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = $1
	`, strings.Join(fieldNames, ","), e.TableName())

	if err := db.QueryRowContext(ctx, stmt, &id).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *productRepository) Create(ctx context.Context, db database.Executor, data *entity.Product) (int64, error) {
	fieldNames, values := database.FieldMap(data)
	fieldNames = fieldNames[1:]
	values = values[1:]
	placeHolders := database.GetPlaceholders(len(fieldNames))

	stmt := fmt.Sprintf(`
		INSERT INTO %s(%s)
		VALUES(%s)
		RETURNING id
	`, data.TableName(), strings.Join(fieldNames, ","), placeHolders)
	var id int64

	if err := db.QueryRowContext(ctx, stmt, values...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *productRepository) UpdateByID(ctx context.Context, db database.Executor, id int64, data *entity.Product) error {
	e := &entity.Product{}
	stmt := fmt.Sprintf(`
		UPDATE %s
		SET
		name = COALESCE($2, name),
		type = COALESCE($3, type),
		description = COALESCE($4, description),
		price = COALESCE($5, price)
		WHERE id = $1
	`, e.TableName())

	result, err := db.ExecContext(ctx, stmt, &id, &data.Name, &data.Type, &data.Description, &data.Price)
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

func (r *productRepository) DeleteByID(ctx context.Context, db database.Executor, id int64) error {
	e := &entity.Product{}
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

func (r *productRepository) DeleteByIDs(ctx context.Context, db database.Executor, ids []int64) error {
	e := &entity.Product{}
	stmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = ANY($1)
	`, e.TableName())

	result, err := db.ExecContext(ctx, stmt, pq.Int64Array(ids))
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

func (r *productRepository) Count(ctx context.Context, db database.Executor) (int64, error) {
	e := &entity.Product{}
	stmt := fmt.Sprintf(`
		SELECT COUNT(1)
		FROM %s
	`, e.TableName())
	var total sql.NullInt64
	if err := db.QueryRowContext(ctx, stmt).Scan(&total); err != nil {
		return total.Int64, err
	}

	return total.Int64, nil
}
