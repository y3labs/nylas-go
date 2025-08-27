package nylas

import (
	"net/http"
	"testing"
	"time"
)

func TestParseRateLimit_NilHeader(t *testing.T) {
	out := ParseRateLimit(nil)
	if out.Limit != nil || out.Remaining != nil || out.ResetAt != nil {
		t.Fatalf("expected zero-value output for nil header, got %#v", out)
	}
}

func TestParseRateLimit_EmptyValues_Continue(t *testing.T) {
	h := make(http.Header)
	// Force len(vs) == 0 to hit the `continue` branch
	h["X-RateLimit-Limit"] = []string{}
	// Provide another header that *does* parse so we know processing continues
	h.Set("X-RateLimit-Remaining", "7")

	out := ParseRateLimit(h)
	if out.Limit != nil {
		t.Fatalf("Limit should be nil when header value slice is empty, got %v", *out.Limit)
	}
	if out.Remaining == nil || *out.Remaining != 7 {
		t.Fatalf("Remaining parse failed, got %#v", out.Remaining)
	}
}

func TestParseRateLimit_NormalHeaders_AllSwitchCases(t *testing.T) {
	h := make(http.Header)
	h.Set("X-RateLimit-Limit", "100")        // case: ratelimit-limit
	h.Set("X-RateLimit-Remaining", "42")     // case: ratelimit-remaining
	h.Set("X-RateLimit-Reset", "1730000000") // case: exact match (seconds epoch)

	out := ParseRateLimit(h)
	if out.Limit == nil || *out.Limit != 100 {
		t.Fatalf("Limit parse failed: %#v", out.Limit)
	}
	if out.Remaining == nil || *out.Remaining != 42 {
		t.Fatalf("Remaining parse failed: %#v", out.Remaining)
	}
	want := time.Unix(1730000000, 0).UTC()
	if out.ResetAt == nil || !out.ResetAt.Equal(want) {
		t.Fatalf("ResetAt parse failed: got %v want %v", out.ResetAt, want)
	}
}

func TestParseRateLimit_CasingVariantsAndRetryAfterFallback(t *testing.T) {
	h := make(http.Header)
	// Variants that your code actually matches (no extra hyphen between "rate" and "limit"):
	h.Set("x-ratelimit-limit", "5")     // contains "ratelimit-limit"
	h.Set("X-RateLimit-Remaining", "3") // contains "ratelimit-remaining"
	// No reset header → trigger Retry-After fallback
	h.Set("Retry-After", "2")

	start := time.Now().UTC()
	out := ParseRateLimit(h)
	end := time.Now().UTC()

	if out.Limit == nil || *out.Limit != 5 {
		t.Fatalf("Limit parse failed with variant header: %#v", out.Limit)
	}
	if out.Remaining == nil || *out.Remaining != 3 {
		t.Fatalf("Remaining parse failed with variant header: %#v", out.Remaining)
	}

	if out.ResetAt == nil {
		t.Fatalf("ResetAt should be set from Retry-After")
	}
	min := start.Add(2 * time.Second).Add(-200 * time.Millisecond)
	max := end.Add(3 * time.Second)
	if out.ResetAt.Before(min) || out.ResetAt.After(max) {
		t.Fatalf("ResetAt from Retry-After out of expected range: got %v, want between %v and %v", *out.ResetAt, min, max)
	}
}
