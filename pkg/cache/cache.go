package cache

import (
	"context"
)

// Cache is an exporter for common interface of memories caching [trintech/review/lru.lru]
// or third party caching like redis.
type Cache[K comparable, V any] interface {
	Add(context.Context, K, V) error
	Get(context.Context, K) (V, error)
	Remove(context.Context, K) error
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}
