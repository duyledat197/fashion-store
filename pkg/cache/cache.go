package cache

import (
	"context"
)

// Cache is an interface that defines common caching operations, such as Add, Get, and Remove.
// It serves as an abstraction for different caching implementations, including in-memory caching (e.g., lru.LRU)
// and third-party caching systems like Redis.
type Cache[K comparable, V any] interface {
	Add(context.Context, K, V) error   // Add adds a key-value pair to the cache.
	Get(context.Context, K) (V, error) // Get retrieves the value associated with the given key from the cache.
	Remove(context.Context, K) error   // Remove removes the key and its associated value from the cache.
}

// number is an interface used to specify a set of numeric types that can be used as key types in the Cache.
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}
