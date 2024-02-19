package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
)

// UserHistoryRepository ...
type UserHistoryRepository interface {
	Create(ctx context.Context, data *entity.UserHistory) (int64, error)
}
