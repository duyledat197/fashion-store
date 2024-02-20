package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

// UserRepository ...
type UserRepository interface {
	RetrieveByEmail(ctx context.Context, db database.Executor, email string) (*entity.User, error)
	RetrieveByUserName(ctx context.Context, db database.Executor, userName string) (*entity.User, error)
	Create(ctx context.Context, db database.Executor, data *entity.User) (int64, error)
	UpdatePassword(ctx context.Context, db database.Executor, email, password string) error
}

// UserCacheRepository ...
type UserCacheRepository interface {
	RetrieveByUserName(context.Context, string) (*entity.User, error)
	StoreByUserName(context.Context, string, *entity.User) error
	RemoveByUserName(context.Context, string) error
	RetrieveByEmail(context.Context, string) (*entity.User, error)
	StoreByEmail(context.Context, string, *entity.User) error
	RemoveByEmail(context.Context, string) error
	IncrementForgotPassword(ctx context.Context, email string) (int64, error)
	StoreResetToken(ctx context.Context, email string, resetToken string) error
	IsExistResetToken(ctx context.Context, email string, resetToken string) error
	RemoveByResetToken(context.Context, string, string) error
}
