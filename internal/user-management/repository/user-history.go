// Package repository defines interfaces related to user management and database operations.
package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
)

// UserHistoryRepository defines methods for creating user history records.
type UserHistoryRepository interface {
	// Create adds a new user history record to the database.
	// It returns the ID of the newly created record and an error if any.
	Create(ctx context.Context, data *entity.UserHistory) (int64, error)
}
