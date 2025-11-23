package middleware

import (
	"net/http"
	"github.com/jerkeyray/tokbuk/internal/limiter"
)

type KeyFunc func(*http.Request) string

func RateLimit(
	manager *limiter.BucketManager,
	capacity int64,
	rate float64,
	keyFunc KeyFunc,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				key := keyFunc(r)
				bucket := manager.GetBucket(key, capacity, rate)

				if !bucket.Allow(1) {
					http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
					return
				}

				next.ServeHTTP(w, r)
			})
		}
}

