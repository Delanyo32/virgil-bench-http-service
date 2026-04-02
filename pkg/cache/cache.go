package cache

import (
	"time"
)

// Cache defines the interface for caching operations.
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Cleanup()
}

// CacheEntry represents a cached item with expiration.
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired returns true if the cache entry has passed its TTL.
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}
