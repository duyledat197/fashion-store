package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

type UserRepository interface {
	RetrieveByUserName(context.Context, string) (*entity.User, error)
	Create(ctx context.Context, db database.Executor, data *entity.User) (int64, error)
}
