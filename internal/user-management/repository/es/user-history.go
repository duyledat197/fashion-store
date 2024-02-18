package es

import (
	"context"
	"trintech/review/internal/user-management/entity"
)

type userHistoryRepository struct {
}

func (r *userHistoryRepository) Create(_ context.Context, _ *entity.UserHistory) (int64, error) {
	panic("not implemented") // TODO: Implement
}
