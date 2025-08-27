package nylas

import (
	"net/http"
	"testing"
	"time"
)

func TestShouldRetry(t *testing.T) {
	c := NewClient("test-key", WithTimeout(50*time.Millisecond))

	// Common retryables
	if !c.ShouldRetry(429) { // allow either receiver or pkg-level depending on your impl
		t.Fatalf("429 should be retryable")
	}
	if c.ShouldRetry(500) {
		t.Fatalf("500 shouldn't be retryable")
	}

	// Non-retryables
	if c.ShouldRetry(400) {
		t.Fatalf("400 should NOT be retryable")
	}
	if c.ShouldRetry(418) {
		t.Fatalf("418 should NOT be retryable")
	}
}

func TestBackoffSleep_UsesRetryAfterZero(t *testing.T) {
	// We don't want tests to sleep long; Retry-After: 0 should return near-immediately.
	h := make(http.Header)
	h.Set("Retry-After", "0")

	c := NewClient("test-key", WithTimeout(50*time.Millisecond))
	start := time.Now()
	c.backoffSleep(0, h, DefaultRetryPolicy())
	elapsed := time.Since(start)

	if elapsed > 50*time.Millisecond {
		t.Fatalf("backoffSleep slept too long with Retry-After=0: %v", elapsed)
	}
}

func TestBackoffSleep_NoRetryAfter_DoesNotExplode(t *testing.T) {
	// Smoke test: ensure it doesn't panic and doesn't sleep excessively (assumes reasonable base delay).
	c := NewClient("test-key", WithTimeout(50*time.Millisecond))
	start := time.Now()
	c.backoffSleep(0, http.Header{}, DefaultRetryPolicy())
	elapsed := time.Since(start)

	// Allow up to 1s to accommodate larger base delays/jitter.
	if elapsed > time.Second {
		t.Fatalf("backoffSleep took too long without Retry-After: %v", elapsed)
	}
}
