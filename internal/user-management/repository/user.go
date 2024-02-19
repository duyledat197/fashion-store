package repository

import (
	"context"
	"time"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

// UserRepository ...
type UserRepository interface {
	RetrieveByEmail(context.Context, string) (*entity.User, error)
	RetrieveByUserName(context.Context, string) (*entity.User, error)
	Create(context.Context, database.Executor, *entity.User) (int64, error)
	UpdatePassword(ctx context.Context, email, password string) error
}

// UserCacheRepository ...
type UserCacheRepository interface {
	RetrieveByUserName(context.Context, string) (*entity.User, error)
	StoreByUserName(context.Context, string, *entity.User) error
	RemoveByUserName(context.Context, string) error
	RetrieveByEmail(context.Context, string) (*entity.User, error)
	StoreByEmail(context.Context, string, *entity.User) error
	RemoveByEmail(context.Context, string) error
	IncrementForgotPassword(ctx context.Context, email string, duration time.Duration) (int64, error)
	StoreResetToken(ctx context.Context, email string, resetToken string, duration time.Duration) error
	RetrieveResetToken(ctx context.Context, email string, resetToken string) error
	RemoveByResetToken(context.Context, string, string) error
}
