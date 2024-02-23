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

// userRepository is an implementation of the UserRepository interface for PostgreSQL database.
type userRepository struct {
}

// NewUserRepository creates a new instance of userRepository.
func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

// RetrieveByUserName retrieves a user from the database based on the username.
// It returns the retrieved user and an error if any.
func (r *userRepository) RetrieveByUserName(ctx context.Context, db database.Executor, userName string) (*entity.User, error) {
	e := &entity.User{}
	fieldNames, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE user_name = $1
	`, strings.Join(fieldNames, ","), e.TableName())

	if err := db.QueryRowContext(ctx, stmt, &userName).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}

// Create adds a new user record to the database.
// It returns the ID of the newly created record and an error if any.
func (r *userRepository) Create(ctx context.Context, db database.Executor, data *entity.User) (int64, error) {
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

// RetrieveByEmail retrieves a user from the database based on the email.
// It returns the retrieved user and an error if any.
func (r *userRepository) RetrieveByEmail(ctx context.Context, db database.Executor, email string) (*entity.User, error) {
	e := &entity.User{}
	fieldNames, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE email = $1
	`, strings.Join(fieldNames, ","), e.TableName())

	if err := db.QueryRowContext(ctx, stmt, &email).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}

// UpdatePassword updates the password of a user in the database based on the email.
// It returns an error if any.
func (r *userRepository) UpdatePassword(ctx context.Context, db database.Executor, email string, password string) error {
	e := &entity.User{}
	stmt := fmt.Sprintf(`
		UPDATE %s
		SET
		password = $2
		WHERE email = $1
	`, e.TableName())
	result, err := db.ExecContext(ctx, stmt, &email, &password)
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
