package api

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter per client IP.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int
	interval time.Duration
}

// bucket tracks tokens for a single client.
type bucket struct {
	tokens    int
	lastFill  time.Time
}

// NewRateLimiter creates a rate limiter allowing the given number of
// requests per interval.
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		interval: interval,
	}
}

// Allow checks whether the client identified by key is within the rate limit.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.buckets[key]
	now := time.Now()

	if !exists {
		rl.buckets[key] = &bucket{tokens: rl.rate - 1, lastFill: now}
		return true
	}

	elapsed := now.Sub(b.lastFill)
	if elapsed >= rl.interval {
		b.tokens = rl.rate
		b.lastFill = now
	}

	if b.tokens <= 0 {
		return false
	}

	b.tokens--
	return true
}

// Middleware returns an HTTP middleware that enforces the rate limit.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		if !rl.Allow(key) {
			WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Reset clears all tracked buckets.
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.buckets = make(map[string]*bucket)
}
