package cache

import (
	"container/list"
	"sync"
	"time"
)

// lruEntry pairs a key with its list element for O(1) eviction.
type lruEntry struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// LRUCache is a thread-safe least-recently-used cache with TTL support.
type LRUCache struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

// NewLRUCache creates an LRU cache with the given maximum capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves a value by key, returning (value, true) on hit.
// Expired entries are treated as misses and evicted.
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}

	entry := elem.Value.(*lruEntry)
	if time.Now().After(entry.expiresAt) {
		c.removeElement(elem)
		return nil, false
	}

	c.order.MoveToFront(elem)
	return entry.value, true
}

// Set adds or updates a cache entry with the given TTL.
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		entry := elem.Value.(*lruEntry)
		entry.value = value
		entry.expiresAt = time.Now().Add(ttl)
		return
	}

	entry := &lruEntry{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	elem := c.order.PushFront(entry)
	c.items[key] = elem

	if c.order.Len() > c.capacity {
		c.evictOldest()
	}
}

// Delete removes an entry by key.
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
	}
}

// Cleanup removes all expired entries.
func (c *LRUCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, elem := range c.items {
		entry := elem.Value.(*lruEntry)
		if now.After(entry.expiresAt) {
			c.order.Remove(elem)
			delete(c.items, key)
		}
	}
}

// Len returns the current number of entries.
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// evictOldest removes the least-recently-used entry.
func (c *LRUCache) evictOldest() {
	elem := c.order.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

// removeElement removes a list element and its map entry.
func (c *LRUCache) removeElement(elem *list.Element) {
	c.order.Remove(elem)
	entry := elem.Value.(*lruEntry)
	delete(c.items, entry.key)
}
