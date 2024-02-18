package lru

import (
	"context"
	"fmt"
	"time"
	"trintech/review/pkg/cache"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

// lru is presentation of implementing lru memories cache of [cache.Cache]
type lru[K comparable, V any] struct {
	*expirable.LRU[K, V]
}

func NewLRU[K comparable, V any](size int, ttl time.Duration) cache.Cache[K, V] {
	return &lru[K, V]{
		expirable.NewLRU[K, V](size, nil, ttl),
	}
}

// Add is implementation of Add by [lru] in [cache.Cache]
func (c *lru[K, V]) Add(_ context.Context, k K, v V) error {
	if ok := c.LRU.Add(k, v); !ok {
		return fmt.Errorf("unable to add value to lru")
	}

	return nil
}

// Get is implementation of Get by [lru] in [cache.Cache]
func (c *lru[K, V]) Get(_ context.Context, k K) (V, error) {
	v, ok := c.LRU.Get(k)
	if !ok {
		return v, fmt.Errorf("value of %v does not exists", k)
	}

	return v, nil
}

// Remove is implementation of Remove by [lru] in [cache.Cache]
func (c *lru[K, V]) Remove(_ context.Context, k K) error {
	if ok := c.LRU.Remove(k); !ok {
		return fmt.Errorf("unable to remove value of %v from lru", k)
	}

	return nil
}
