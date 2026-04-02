package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery returns an HTTP middleware that recovers from panics,
// logs the stack trace, and returns a 500 Internal Server Error.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Printf("panic recovered: %v\n%s", err, stack)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"internal server error"}`)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RecoveryWithHandler returns a recovery middleware that calls a custom
// error handler instead of writing a default response.
func RecoveryWithHandler(handler func(w http.ResponseWriter, r *http.Request, err interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					log.Printf("panic recovered: %v\n%s", err, stack)
					handler(w, r, err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
