package nylas

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RateLimitInfo struct {
	Limit     *int
	Remaining *int
	ResetAt   *time.Time
}

func ParseRateLimit(h http.Header) RateLimitInfo {
	var out RateLimitInfo
	if h == nil {
		return out
	}

	// Normalize and scan all headers to catch spelling/casing variants:
	// X-RateLimit-Limit, X-Rate-Limit-Limit, x-ratelimit-limit, etc.
	for k, vs := range h {
		if len(vs) == 0 {
			continue
		}
		kl := strings.ToLower(k)
		val := vs[0]

		switch {
		case strings.Contains(kl, "ratelimit-limit"):
			if i, err := strconv.Atoi(val); err == nil {
				out.Limit = &i
			}
		case strings.Contains(kl, "ratelimit-remaining"):
			if i, err := strconv.Atoi(val); err == nil {
				out.Remaining = &i
			}
		case kl == "x-ratelimit-reset" || kl == "x-rate-limit-reset":
			if secs, err := strconv.ParseInt(val, 10, 64); err == nil {
				t := time.Unix(secs, 0).UTC()
				out.ResetAt = &t
			}
		}
	}

	// If ResetAt still empty, fall back to Retry-After seconds from now.
	if out.ResetAt == nil {
		if ra := h.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				t := time.Now().Add(time.Duration(secs) * time.Second).UTC()
				out.ResetAt = &t
			}
		}
	}

	return out
}
