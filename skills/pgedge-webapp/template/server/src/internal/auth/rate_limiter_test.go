package auth

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !rl.Allow("1.2.3.4") {
			t.Fatalf("attempt %d should be allowed", i+1)
		}
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("4th attempt should be rate-limited")
	}
	if !rl.Allow("9.9.9.9") {
		t.Fatal("different IP should not be limited")
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	if !rl.Allow("k") {
		t.Fatal("first allowed")
	}
	if rl.Allow("k") {
		t.Fatal("second blocked")
	}
	time.Sleep(75 * time.Millisecond)
	if !rl.Allow("k") {
		t.Fatal("post-window should be allowed")
	}
}
