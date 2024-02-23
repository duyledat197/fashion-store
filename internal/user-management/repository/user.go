// Package repository defines interfaces related to user management and database operations.
package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

// UserRepository defines methods for interacting with user data in the database.
type UserRepository interface {
	// RetrieveByEmail fetches a user record from the database based on the email.
	// It returns the retrieved user and an error if any.
	RetrieveByEmail(ctx context.Context, db database.Executor, email string) (*entity.User, error)

	// RetrieveByUserName fetches a user record from the database based on the username.
	// It returns the retrieved user and an error if any.
	RetrieveByUserName(ctx context.Context, db database.Executor, userName string) (*entity.User, error)

	// Create adds a new user record to the database.
	// It returns the ID of the newly created record and an error if any.
	Create(ctx context.Context, db database.Executor, data *entity.User) (int64, error)

	// UpdatePassword updates the password of a user in the database based on the email.
	// It returns an error if any.
	UpdatePassword(ctx context.Context, db database.Executor, email, password string) error
}

// UserCacheRepository defines methods for caching and retrieving user-related data.
type UserCacheRepository interface {
	// RetrieveByUserName retrieves a user from the cache based on the username.
	// It returns the retrieved user and an error if any.
	RetrieveByUserName(context.Context, string) (*entity.User, error)

	// StoreByUserName stores a user in the cache based on the username.
	// It returns an error if any.
	StoreByUserName(context.Context, string, *entity.User) error

	// RemoveByUserName removes a user from the cache based on the username.
	// It returns an error if any.
	RemoveByUserName(context.Context, string) error

	// RetrieveByEmail retrieves a user from the cache based on the email.
	// It returns the retrieved user and an error if any.
	RetrieveByEmail(context.Context, string) (*entity.User, error)

	// StoreByEmail stores a user in the cache based on the email.
	// It returns an error if any.
	StoreByEmail(context.Context, string, *entity.User) error

	// RemoveByEmail removes a user from the cache based on the email.
	// It returns an error if any.
	RemoveByEmail(context.Context, string) error

	// IncrementForgotPassword increments the count of forgot password attempts for a user.
	// It returns the updated count and an error if any.
	IncrementForgotPassword(ctx context.Context, email string) (int64, error)

	// StoreResetToken stores a reset token in the cache for a user.
	// It returns an error if any.
	StoreResetToken(ctx context.Context, email string, resetToken string) error

	// IsExistResetToken checks if a reset token exists in the cache for a user.
	// It returns an error if any.
	IsExistResetToken(ctx context.Context, email string, resetToken string) error

	// RemoveByResetToken removes a user from the cache based on the reset token.
	// It returns an error if any.
	RemoveByResetToken(context.Context, string, string) error
}
