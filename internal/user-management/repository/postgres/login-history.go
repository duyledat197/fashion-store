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

// loginHistoryRepository is an implementation of the LoginHistoryRepository interface for PostgreSQL database.
type loginHistoryRepository struct {
}

// NewLoginHistoryRepository creates a new instance of loginHistoryRepository.
func NewLoginHistoryRepository() repository.LoginHistoryRepository {
	return &loginHistoryRepository{}
}

// Create adds a new login history record to the database.
// It returns an error if any.
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

// UpdateLogout updates the logout timestamp of a user's login history in the database based on the access token and user ID.
// It returns an error if any.
func (r *loginHistoryRepository) UpdateLogout(ctx context.Context, db database.Executor, userID int64, accessToken string) error {
	e := &entity.LoginHistory{}
	stmt := fmt.Sprintf(`
		UPDATE %s
		SET
		logout_at = NOW()
		WHERE access_token = $1
		AND user_id = $2
	`, e.TableName())
	result, err := db.ExecContext(ctx, stmt, &accessToken, &userID)
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
