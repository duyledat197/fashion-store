package cache

import "context"

// Cache is an exporter for common interface of memories caching [trintech/review/lru.lru]
// or third party caching like redis.
type Cache[K comparable, V any] interface {
	Add(context.Context, K, V) error
	Get(context.Context, K) (V, error)
	Remove(context.Context, K) error
}
