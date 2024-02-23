// Package repository defines interfaces related to user management and database operations.
package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

// LoginHistoryRepository defines methods for creating and updating login history records.
type LoginHistoryRepository interface {
	// Create adds a new login history record to the database.
	Create(ctx context.Context, db database.Executor, data *entity.LoginHistory) error

	// UpdateLogout updates the logout information in the login history records.
	// It marks the session as logged out for the specified user ID and access token.
	UpdateLogout(ctx context.Context, db database.Executor, userID int64, accessToken string) error
}
