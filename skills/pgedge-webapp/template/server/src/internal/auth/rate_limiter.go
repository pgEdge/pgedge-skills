package auth

import (
	"sync"
	"time"
)

// RateLimiter is an in-memory sliding-window per-key limiter.
type RateLimiter struct {
	max    int
	window time.Duration
	mu     sync.Mutex
	hits   map[string][]time.Time
}

// NewRateLimiter constructs a limiter allowing max attempts per window.
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	return &RateLimiter{max: max, window: window, hits: make(map[string][]time.Time)}
}

// Allow records an attempt for key and returns true if it is within the cap.
func (rl *RateLimiter) Allow(key string) bool {
	if rl.max <= 0 || rl.window <= 0 {
		return true
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-rl.window)
	times := rl.hits[key]
	pruned := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			pruned = append(pruned, t)
		}
	}
	if len(pruned) >= rl.max {
		rl.hits[key] = pruned
		return false
	}
	pruned = append(pruned, now)
	rl.hits[key] = pruned
	return true
}
