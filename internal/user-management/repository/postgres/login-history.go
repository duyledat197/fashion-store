package postgres

import (
	"context"

	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

type loginHistoryRepository struct {
}

func (r *loginHistoryRepository) Create(_ context.Context, db database.Executor, _ *entity.LoginHistory) (int64, error) {
	panic("not implemented") // TODO: Implement
}
