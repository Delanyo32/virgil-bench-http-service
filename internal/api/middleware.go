package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/example/ordersvc/pkg/logger"
)

// LoggingMiddleware logs request details.
func LoggingMiddleware(logr *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logr.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration.String(),
				"remote", r.RemoteAddr,
			)
		})
	}
}

// AuthMiddleware validates JWT tokens from the Authorization header.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health and ready endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/ready" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, http.StatusUnauthorized, "missing authorization")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				WriteError(w, http.StatusUnauthorized, "invalid auth format")
				return
			}

			// Simple token validation (not real JWT)
			token := parts[1]
			if len(token) < 10 {
				WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// rateLimitEntry tracks per-client request counts.
type rateLimitEntry struct {
	count    int
	resetAt  time.Time
}

// RateLimitMiddleware limits requests per client IP.
// FLAW: uses a shared map with sync.Mutex but the cleanup goroutine
// can race with request handlers during map iteration.
func RateLimitMiddleware(maxRequests int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]*rateLimitEntry)

	// Cleanup goroutine -- iterates map while handlers may write to it.
	// FLAW: goroutine leak -- no way to stop this goroutine.
	go func() {
		for {
			time.Sleep(60 * time.Second) // FLAW: magic number
			mu.Lock()
			now := time.Now()
			for ip, entry := range clients {
				if now.After(entry.resetAt) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr

			mu.Lock()
			entry, exists := clients[clientIP]
			if !exists {
				entry = &rateLimitEntry{
					count:   0,
					resetAt: time.Now().Add(60 * time.Second),
				}
				clients[clientIP] = entry
			}

			if time.Now().After(entry.resetAt) {
				entry.count = 0
				entry.resetAt = time.Now().Add(60 * time.Second)
			}

			entry.count++
			count := entry.count
			mu.Unlock()

			if count > maxRequests {
				w.Header().Set("Retry-After", "60")
				WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware adds CORS headers. Dead code -- not used in router setup.
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware catches panics and returns 500.
// Dead code -- not used in router setup.
func RecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Printf("panic recovered: %v", rec)
					WriteError(w, http.StatusInternalServerError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adds a request ID header.
// Dead code -- not wired into middleware chain.
func RequestIDMiddleware() func(http.Handler) http.Handler {
	var counter int64
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			counter++
			id := counter
			mu.Unlock()

			r.Header.Set("X-Request-ID", fmt.Sprintf("req-%d", id))
			next.ServeHTTP(w, r)
		})
	}
}
