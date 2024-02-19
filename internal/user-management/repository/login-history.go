package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

// LoginHistoryRepository ...
type LoginHistoryRepository interface {
	Create(ctx context.Context, db database.Executor, data *entity.LoginHistory) error
	UpdateLogout(ctx context.Context, db database.Executor, accessToken string) error
}
