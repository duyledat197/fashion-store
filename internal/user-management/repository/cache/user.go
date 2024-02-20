package memcache

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"trintech/review/internal/user-management/entity"
	"trintech/review/internal/user-management/repository"
	"trintech/review/pkg/cache"
	"trintech/review/pkg/lru"
)

type userCacheRepository struct {
	cache cache.Cache[string, *entity.User]

	fpMap cache.Cache[string, *int64] // forgot password map
	rsMap cache.Cache[string, bool]   // reset token map
}

// NewUserCacheRepository ...
func NewUserCacheRepository() repository.UserCacheRepository {
	return &userCacheRepository{
		cache: lru.NewLRU[string, *entity.User](1000, 10*time.Minute),
		fpMap: lru.NewLRU[string, *int64](1000, 5*time.Minute),
		rsMap: lru.NewLRU[string, bool](1000, 5*time.Minute),
	}
}
func (r *userCacheRepository) RetrieveByUserName(ctx context.Context, userName string) (*entity.User, error) {
	user, err := r.cache.Get(ctx, fmt.Sprintf("userName|%s", userName))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userCacheRepository) StoreByUserName(ctx context.Context, userName string, user *entity.User) error {
	if err := r.cache.Add(ctx, fmt.Sprintf("userName|%s", userName), user); err != nil {
		return err
	}

	return nil
}

func (r *userCacheRepository) RemoveByUserName(ctx context.Context, userName string) error {
	if err := r.cache.Remove(ctx, fmt.Sprintf("userName|%s", userName)); err != nil {
		return err
	}

	return nil
}

func (r *userCacheRepository) RetrieveByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.cache.Get(ctx, fmt.Sprintf("email|%s", email))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userCacheRepository) StoreByEmail(ctx context.Context, email string, user *entity.User) error {
	if err := r.cache.Add(ctx, fmt.Sprintf("email|%s", email), user); err != nil {
		return err
	}

	return nil
}

func (r *userCacheRepository) RemoveByEmail(ctx context.Context, email string) error {
	if err := r.cache.Remove(ctx, fmt.Sprintf("email|%s", email)); err != nil {
		return err
	}

	return nil
}

func (r *userCacheRepository) IncrementForgotPassword(ctx context.Context, email string) (int64, error) {
	num, err := r.fpMap.Get(ctx, email)
	if err != nil {
		return 0, err
	}

	if num == nil {
		num = new(int64)
	}

	atomic.AddInt64(num, 1)

	return *num, nil
}

func (r *userCacheRepository) StoreResetToken(ctx context.Context, email string, resetToken string) error {
	key := fmt.Sprintf("%s|%s", email, resetToken)

	if err := r.rsMap.Add(ctx, key, true); err != nil {
		return err
	}

	return nil
}

func (r *userCacheRepository) IsExistResetToken(ctx context.Context, email string, resetToken string) error {
	key := fmt.Sprintf("%s|%s", email, resetToken)
	isExist, err := r.rsMap.Get(ctx, key)
	if err != nil {
		return err
	}

	if !isExist {
		return fmt.Errorf("reset token not exists")
	}

	return nil
}

func (r *userCacheRepository) RemoveByResetToken(ctx context.Context, email string, resetToken string) error {
	key := fmt.Sprintf("%s|%s", email, resetToken)

	if err := r.rsMap.Remove(ctx, key); err != nil {
		return err
	}

	return nil
}
