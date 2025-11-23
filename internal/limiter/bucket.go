package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu       sync.Mutex       // allow only one goroutine to read/modify this bucket at a time
	capacity int64            // burst allowance (max tokens bucket can hold)
	tokens   float64          // current available tokens (float because refill can be fractional)
	rate     float64          // tokens regenerated per second
	last     time.Time        // when last refill was calculated
	now      func() time.Time // injected time source so tests can control time
}

// create bucket, fill it to max capacity, record time of creation and set refill rate
func NewTokenBucket(capacity int64, ratePerSecond float64) *TokenBucket {
	if capacity <= 0 {
		panic("capacity must be greater than 0!")
	}
	if ratePerSecond <= 0 {
		panic("ratePerSecond must be greater than 0!")
	}

	now := time.Now() // used to calculate elapsed time for refill

	return &TokenBucket{
		capacity: capacity,
		tokens:   float64(capacity),
		rate:     ratePerSecond,
		last:     now,
		now:      time.Now,
	}
}

// return true if request is permitted else false
func (b *TokenBucket) Allow(n int64) bool {
	if n <= 0 {
		return false // reject zero or negative token requests
	}

	// lock the mutex
	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill()

	// if enough tokens available spend them and return true
	if float64(n) <= b.tokens {
		b.tokens -= float64(n)
		return true
	}

	// deny request
	return false
}

// generate tokens based on time passed since last refill
func (b *TokenBucket) refill() {
	now := b.now()                       // current time
	elapsed := now.Sub(b.last).Seconds() // current time - timestamp of last refill
	if elapsed <= 0 {
		return
	}

	add := elapsed * b.rate // tokens regenerated based on time passed
	if add > 0 {
		b.tokens += add // add them to the bucket
		// if bucket overflows set tokens to max capacity
		if b.tokens > float64(b.capacity) {
			b.tokens = float64(b.capacity)
		}
		b.last = now // update last refill timestamp
	}
}

// replace real time with a controllable clock for deterministic testing
func (b *TokenBucket) WithClock(now func() time.Time) *TokenBucket {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.now = now    // swap out the time provider
	b.last = now() // resets last
	return b
}
