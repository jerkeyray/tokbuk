package limiter 

import (
	"sync"
)

type BucketManager struct {
	mu sync.Mutex
	buckets map[string]*TokenBucket
}

func NewBucketManager() *BucketManager {
	return &BucketManager {
		buckets: make(map[string]*TokenBucket),
	}
}

func (m *BucketManager) GetBucket(key string, capacity int64, rate float64) *TokenBucket {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if bucket, ok := m.buckets[key]; ok {
		return bucket
	}
	
	bucket := NewTokenBucket(capacity, rate)
	m.buckets[key] = bucket
	return bucket
}