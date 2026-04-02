package cache

import (
	"log"
	"time"
)

// RedisCache implements Cache using an in-memory map.
// Despite the name, this does not actually use Redis -- it is a
// naive in-memory implementation that was intended to be replaced.
type RedisCache struct {
	addr      string
	maxSize   int
	// FLAW: race condition -- this map is accessed concurrently by HTTP
	// handlers (Get/Set) and the cleanup goroutine (Cleanup) without
	// any synchronization. This causes concurrent map read/write panics
	// under load.
	data map[string]*CacheEntry
}

// NewRedisCache creates a new RedisCache with the given address and max size.
func NewRedisCache(addr string, maxSize int) *RedisCache {
	return &RedisCache{
		addr:    addr,
		maxSize: maxSize,
		data:    make(map[string]*CacheEntry),
	}
}

// Get retrieves a value from the cache.
// FLAW: concurrent map read without mutex -- panics under concurrent access.
func (c *RedisCache) Get(key string) (interface{}, bool) {
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	if entry.IsExpired() {
		delete(c.data, key)
		return nil, false
	}
	return entry.Value, true
}

// Set stores a value in the cache with a TTL.
// FLAW: concurrent map write without mutex -- panics under concurrent access.
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) {
	// FLAW: no size limit enforcement -- maxSize is stored but never checked.
	// Cache grows unbounded in memory.
	c.data[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a key from the cache.
func (c *RedisCache) Delete(key string) {
	delete(c.data, key)
}

// Cleanup removes expired entries from the cache.
// FLAW: iterates map while other goroutines may be reading/writing,
// causing concurrent map iteration panics.
func (c *RedisCache) Cleanup() {
	count := 0
	for key, entry := range c.data {
		if entry.IsExpired() {
			delete(c.data, key)
			count++
		}
	}
	if count > 0 {
		log.Printf("cache cleanup: removed %d expired entries", count)
	}
}

// Size returns the number of entries in the cache.
func (c *RedisCache) Size() int {
	return len(c.data)
}

// Flush removes all entries from the cache.
// Dead code -- not called from any module.
func (c *RedisCache) Flush() {
	c.data = make(map[string]*CacheEntry)
	log.Println("cache flushed")
}
