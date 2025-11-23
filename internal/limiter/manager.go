package limiter

import (
	"sync"
)

// maintain a collection of token buckets
// one per unique client (IP, userID, APIKey, etc)
// ensure each client gets its own unique rate limiter instead of sharing a global one
type BucketManager struct {
	mu      sync.Mutex
	buckets map[string]*TokenBucket // map client key -> its token bucket
}

func NewBucketManager() *BucketManager {
	return &BucketManager{
		buckets: make(map[string]*TokenBucket),
	}
}

func (m *BucketManager) GetBucket(key string, capacity int64, rate float64) *TokenBucket {
	m.mu.Lock() // prevent concurrent access to the same map
	defer m.mu.Unlock()

	// return existing bucket if key already exists
	if bucket, ok := m.buckets[key]; ok {
		return bucket
	}

	// if key doesn't exist, create and register a new bucket
	// newly created buckets start full and use the provided capacity/rate
	bucket := NewTokenBucket(capacity, rate)
	m.buckets[key] = bucket
	return bucket
}
