package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

type LoginHistoryRepo interface {
	Create(ctx context.Context, db database.Executor, data *entity.LoginHistory) (int64, error)
}
