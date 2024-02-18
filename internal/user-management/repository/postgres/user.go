package postgres

import (
	"context"
	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
)

type userRepository struct {
}

func (r *userRepository) RetrieveByUserName(_ context.Context, _ string) (*entity.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *userRepository) Create(ctx context.Context, db database.Executor, data *entity.User) (int64, error) {
	panic("not implemented") // TODO: Implement
}
