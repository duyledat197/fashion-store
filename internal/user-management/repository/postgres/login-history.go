package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"trintech/review/internal/user-management/entity"
	"trintech/review/internal/user-management/repository"
	"trintech/review/pkg/database"
)

type loginHistoryRepository struct {
}

func NewLoginHistoryRepository() repository.LoginHistoryRepository {
	return &loginHistoryRepository{}
}

func (r *loginHistoryRepository) Create(ctx context.Context, db database.Executor, data *entity.LoginHistory) error {
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

func (r *loginHistoryRepository) UpdateLogout(ctx context.Context, db database.Executor, accessToken string) error {
	e := &entity.LoginHistory{}
	stmt := fmt.Sprintf(`
		UPDATE %s
		SET
		logout_at = NOW()
		WHERE access_token = $1
	`, e.TableName())
	result, err := db.ExecContext(ctx, stmt, &accessToken)
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
