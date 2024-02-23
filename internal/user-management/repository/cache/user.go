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

// userCacheRepository is an implementation of the repository.UserCacheRepository interface.
type userCacheRepository struct {
	cache cache.Cache[string, *entity.User] // Cache for storing user information
	fpMap cache.Cache[string, *int64]       // Cache for storing forgot password attempts
	rsMap cache.Cache[string, bool]         // Cache for storing reset tokens
}

// NewUserCacheRepository creates a new instance of userCacheRepository.
func NewUserCacheRepository() repository.UserCacheRepository {
	return &userCacheRepository{
		cache: lru.NewLRU[string, *entity.User](1000, 10*time.Minute),
		fpMap: lru.NewLRU[string, *int64](1000, 5*time.Minute),
		rsMap: lru.NewLRU[string, bool](1000, 5*time.Minute),
	}
}

// RetrieveByUserName retrieves user information from the cache based on the username.
func (r *userCacheRepository) RetrieveByUserName(ctx context.Context, userName string) (*entity.User, error) {
	user, err := r.cache.Get(ctx, fmt.Sprintf("userName|%s", userName))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// StoreByUserName stores user information in the cache based on the username.
func (r *userCacheRepository) StoreByUserName(ctx context.Context, userName string, user *entity.User) error {
	if err := r.cache.Add(ctx, fmt.Sprintf("userName|%s", userName), user); err != nil {
		return err
	}

	return nil
}

// RemoveByUserName removes user information from the cache based on the username.
func (r *userCacheRepository) RemoveByUserName(ctx context.Context, userName string) error {
	if err := r.cache.Remove(ctx, fmt.Sprintf("userName|%s", userName)); err != nil {
		return err
	}

	return nil
}

// RetrieveByEmail retrieves user information from the cache based on the email.
func (r *userCacheRepository) RetrieveByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.cache.Get(ctx, fmt.Sprintf("email|%s", email))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// StoreByEmail stores user information in the cache based on the email.
func (r *userCacheRepository) StoreByEmail(ctx context.Context, email string, user *entity.User) error {
	if err := r.cache.Add(ctx, fmt.Sprintf("email|%s", email), user); err != nil {
		return err
	}

	return nil
}

// RemoveByEmail removes user information from the cache based on the email.
func (r *userCacheRepository) RemoveByEmail(ctx context.Context, email string) error {
	if err := r.cache.Remove(ctx, fmt.Sprintf("email|%s", email)); err != nil {
		return err
	}

	return nil
}

// IncrementForgotPassword increments the count of forgot password attempts for a given email.
func (r *userCacheRepository) IncrementForgotPassword(ctx context.Context, email string) (int64, error) {
	num, _ := r.fpMap.Get(ctx, email)

	if num == nil {
		num = new(int64)
	}

	atomic.AddInt64(num, 1)

	return *num, nil
}

// StoreResetToken stores a reset token in the cache for a given email.
func (r *userCacheRepository) StoreResetToken(ctx context.Context, email string, resetToken string) error {
	key := fmt.Sprintf("%s|%s", email, resetToken)

	if err := r.rsMap.Add(ctx, key, true); err != nil {
		return err
	}

	return nil
}

// IsExistResetToken checks if a reset token exists in the cache for a given email.
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

// RemoveByResetToken removes a reset token from the cache for a given email.
func (r *userCacheRepository) RemoveByResetToken(ctx context.Context, email string, resetToken string) error {
	key := fmt.Sprintf("%s|%s", email, resetToken)

	if err := r.rsMap.Remove(ctx, key); err != nil {
		return err
	}

	return nil
}
