package nylas

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type RetryPolicy struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{MaxRetries: 3, BaseDelay: 300 * time.Millisecond, MaxDelay: 4 * time.Second}
}

func (c *Client) ShouldRetry(status int) bool {
	switch status {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (c *Client) backoffSleep(attempt int, headers http.Header, policy *RetryPolicy) {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}
	if ra := headers.Get("Retry-After"); ra != "" {
		if secs, err := strconv.Atoi(ra); err == nil {
			time.Sleep(time.Duration(secs) * time.Second)
			return
		}
	}
	d := float64(policy.BaseDelay) * math.Pow(2, float64(attempt))
	jitter := 0.5 + rand.Float64()*0.5
	dur := time.Duration(d * jitter)
	if dur > policy.MaxDelay {
		dur = policy.MaxDelay
	}
	time.Sleep(dur)
}
