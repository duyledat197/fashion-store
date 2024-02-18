package repository

import (
	"context"

	"trintech/review/internal/user-management/entity"
)

type UserHistoryRepo interface {
	Create(ctx context.Context, data *entity.UserHistory) (int64, error)
}
