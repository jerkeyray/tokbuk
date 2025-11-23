package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu sync.Mutex
	capacity int64
	tokens float64
	rate float64
	last time.Time
	now		func() time.Time
}

func NewTokenBucket(capacity int64, ratePerSecond float64) *TokenBucket {
	if capacity <= 0 {
		panic("capacity must be greater than 0!")
	}
	if ratePerSecond <= 0 {
		panic("ratePerSecond must be greater than 0!")
	}
	
	now := time.Now()
	
	return &TokenBucket {
		capacity: capacity, 
		tokens: float64(capacity),
		rate: ratePerSecond,
		last: now,
		now: time.Now,
	}
}

func (b *TokenBucket) Allow(n int64) bool {
	if n <= 0 {
		return false
	}
	
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.refill()
	
	if float64(n) <= b.tokens {
		b.tokens -= float64(n)
		return true
	}
	
	return false
}

func (b *TokenBucket) refill() {
	now := b.now()
	elapsed := now.Sub(b.last).Seconds()
	if elapsed <= 0 {
		return 
	}
	
	add := elapsed * b.rate
	if add > 0 {
		b.tokens += add
		if b.tokens > float64(b.capacity) {
			b.tokens = float64(b.capacity)
		}
		b.last = now
	}
}

func (b *TokenBucket) WithClock(now func() time.Time) *TokenBucket {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.now = now
	b.last = now()
	return b
}