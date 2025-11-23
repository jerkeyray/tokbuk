package limiter

import (
	"testing"
	"time"
)

type fakeClock struct {
	t time.Time
}

func (c *fakeClock) Now() time.Time {
	return c.t
}

func (c *fakeClock) Add(d time.Duration) {
	c.t = c.t.Add(d)
}

func newTestBucket(capacity int64, rate float64, start time.Time) (*TokenBucket, *fakeClock) {
	clock := &fakeClock{t: start}
	b := NewTokenBucket(capacity, rate)
	b.WithClock(clock.Now)
	return b, clock
}

func TestInitialCapacityFull(t *testing.T) {
	start := time.Unix(0, 0)
	b, _ := newTestBucket(10, 5, start)

	for i := 0; i < 10; i++ {
		if !b.Allow(1) {
			t.Fatalf("expected request %d to be allowed", i)
		}
	}

	if b.Allow(1) {
		t.Fatal("expected request to be rejected when bucket is empty")
	}
}

func TestRefillOverTime(t *testing.T) {
	start := time.Unix(0, 0)
	b, clock := newTestBucket(10, 2, start)

	for i := 0; i < 10; i++ {
		if !b.Allow(1) {
			t.Fatalf("expected initial token %d to be allowed", i)
		}
	}

	if b.Allow(1) {
		t.Fatal("expected no tokens left")
	}

	clock.Add(1 * time.Second)

	if !b.Allow(1) {
		t.Fatal("expected 1 token after 1 second")
	}

	clock.Add(1 * time.Second)

	if !b.Allow(1) {
		t.Fatal("expected another token after 2 seconds total")
	}
}

func TestDoesNotExceedCapacity(t *testing.T) {
	start := time.Unix(0, 0)
	b, clock := newTestBucket(5, 10, start)

	for i := 0; i < 5; i++ {
		if !b.Allow(1) {
			t.Fatalf("expected initial token %d to be allowed", i)
		}
	}

	if b.Allow(1) {
		t.Fatal("expected to be empty")
	}

	clock.Add(10 * time.Second)

	for i := 0; i < 5; i++ {
		if !b.Allow(1) {
			t.Fatalf("expected token %d after refill", i)
		}
	}

	if b.Allow(1) {
		t.Fatal("expected exactly capacity, not more")
	}
}

func TestAllowMoreThanCapacityFails(t *testing.T) {
	start := time.Unix(0, 0)
	b, _ := newTestBucket(5, 1, start)

	if b.Allow(6) {
		t.Fatal("should not allow request larger than capacity")
	}
}

func TestZeroOrNegativeRequestsAreRejected(t *testing.T) {
	start := time.Unix(0, 0)
	b, _ := newTestBucket(5, 1, start)

	if b.Allow(0) {
		t.Fatal("should reject zero token request")
	}

	if b.Allow(-1) {
		t.Fatal("should reject negative token request")
	}
}