package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jerkeyray/tokbuk/internal/limiter"
)

func TestRateLimitMiddleware(t *testing.T) {
	manager := limiter.NewBucketManager()

	// tiny bucket so we hit limit fast
	capacity := int64(2)
	rate := float64(1)

	// static key so all requests hit same bucket
	keyFunc := func(r *http.Request) string {
		return "test-user"
	}

	rl := RateLimit(manager, capacity, rate, keyFunc)

	// handler that always succeeds if middleware allows it
	handler := rl(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// helper to perform a request
	doReq := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		return rr
	}

	// first two requests should pass
	if rr := doReq(); rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if rr := doReq(); rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	// next one should fail
	if rr := doReq(); rr.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rr.Code)
	}

	// advance time to allow refill
	// NOTE: we need access to bucket to manipulate clock
	bucket := manager.GetBucket("test-user", capacity, rate)

	fakeNow := time.Now()
	bucket.WithClock(func() time.Time { return fakeNow })

	// advance 1 second to refill one token
	fakeNow = fakeNow.Add(1 * time.Second)

	// now it should pass again
	if rr := doReq(); rr.Code != http.StatusOK {
		t.Fatalf("expected 200 after refill, got %d", rr.Code)
	}
}
