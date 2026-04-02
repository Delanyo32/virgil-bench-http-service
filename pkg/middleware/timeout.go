package middleware

import (
	"context"
	"net/http"
	"time"
)

// Timeout returns an HTTP middleware that enforces a request deadline.
// If the handler does not complete within the given duration, the
// request context is cancelled and a 503 Service Unavailable is returned.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			done := make(chan struct{})
			tw := &timeoutWriter{
				ResponseWriter: w,
				headerWritten:  false,
			}

			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				if !tw.headerWritten {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte(`{"error":"request timeout"}`))
				}
			}
		})
	}
}

// timeoutWriter wraps http.ResponseWriter to track whether headers
// have been written before the timeout fires.
type timeoutWriter struct {
	http.ResponseWriter
	headerWritten bool
}

// WriteHeader records that headers have been sent.
func (tw *timeoutWriter) WriteHeader(code int) {
	tw.headerWritten = true
	tw.ResponseWriter.WriteHeader(code)
}

// Write records that response body is being written.
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.headerWritten = true
	return tw.ResponseWriter.Write(b)
}
