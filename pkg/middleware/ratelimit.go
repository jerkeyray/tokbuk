package middleware

import (
	"github.com/jerkeyray/tokbuk/internal/limiter"
	"net/http"
)

// extracts client's unique identity string from the HTTP request
// the identity determines which token bucket will be used by the request
type KeyFunc func(*http.Request) string

// returns middleware that applies rate limiting to an http.Handler
func RateLimit(
	manager *limiter.BucketManager,
	capacity int64,
	rate float64,
	keyFunc KeyFunc,
) func(http.Handler) http.Handler {

	// return middleware wrapper
	return func(next http.Handler) http.Handler {
		// return handler that enforces rate limiting
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)                                // extract identity of caller
			bucket := manager.GetBucket(key, capacity, rate) // get existing bucket or create a new one if new key

			// check if request is allowed
			if !bucket.Allow(1) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// if allowed proceed to actual request handler
			next.ServeHTTP(w, r)
		})
	}
}
