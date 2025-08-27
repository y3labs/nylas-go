package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// --- helpers ---

// --- tests ---

func TestDoStream_BuildURLError(t *testing.T) {
	// Empty server URL makes buildURL fail.
	c := NewClient("test-key", WithServerURL(""))
	_, err := c.doStream(context.Background(), http.MethodGet, "/x", nil, nil)
	if err == nil {
		t.Fatalf("expected error from buildURL")
	}
}

func TestDoStream_NewRequestError(t *testing.T) {
	// Use a valid server but an invalid HTTP method to make http.NewRequest fail.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	_, err := c.doStream(context.Background(), "BAD METHOD", "/x", nil, nil)
	if err == nil {
		t.Fatalf("expected error from http.NewRequest")
	}
}

func TestDoStream_HeadersCopiedAndSuccess(t *testing.T) {
	// Server asserts headers were copied through nested loops, then returns 200.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-A") != "a1" {
			t.Fatalf("missing header X-A")
		}
		// Multi-value header should append both values.
		got := r.Header.Values("X-Multi")
		if len(got) != 2 || got[0] != "v1" || got[1] != "v2" {
			t.Fatalf("multi header not propagated: %#v", got)
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	h := make(http.Header)
	h.Add("X-A", "a1")
	h.Add("X-Multi", "v1")
	h.Add("X-Multi", "v2")

	resp, err := c.doStream(context.Background(), http.MethodGet, "/ok", nil, h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestDoStream_TransportError_WrappedAndRetries(t *testing.T) {
	// RoundTripper always returns a timeout error. doStream should wrap via wrapTransportError
	// and after retries return that wrapped error.
	cl := &http.Client{
		Transport: rtErr{&timeoutErr{}},
		Timeout:   200 * time.Millisecond,
	}
	c := NewClient("test-key", WithServerURL("http://example.invalid"), WithHTTPClient(cl))

	_, err := c.doStream(context.Background(), http.MethodGet, "/x", nil, nil)
	if err == nil {
		t.Fatalf("expected transport error")
	}
	if _, ok := IsTimeoutError(err); !ok {
		t.Fatalf("expected SDKTimeoutError, got %T: %v", err, err)
	}
}

func TestDoStream_StatusRetryBehavior(t *testing.T) {
	pol := DefaultRetryPolicy()

	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if pol.MaxRetries >= 1 {
			// First call: 429 (retryable by virtually any policy)
			if calls == 1 {
				w.WriteHeader(http.StatusTooManyRequests) // 429
				_, _ = w.Write([]byte(`{"error":{"type":"rate_limited","message":"boom"}}`))
				return
			}
			// Second call: succeed
			w.WriteHeader(http.StatusOK)
			return
		}

		// No retries configured → always fail once so we can assert no retry occurs.
		w.WriteHeader(http.StatusTooManyRequests) // 429
		_, _ = w.Write([]byte(`{"error":{"type":"rate_limited","message":"boom"}}`))
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL), WithTimeout(2*time.Second))

	resp, err := c.doStream(context.Background(), http.MethodGet, "/retry", nil, nil)

	if pol.MaxRetries == 0 {
		// Expect immediate error, no retry
		if err == nil {
			t.Fatalf("expected error with MaxRetries=0")
		}
		if calls != 1 {
			t.Fatalf("expected 1 call (no retries), got %d", calls)
		}
		if _, ok := IsAPIError(err); !ok {
			t.Fatalf("expected APIError, got %T", err)
		}
		return
	}

	// With retries: should succeed on the second call.
	if err != nil {
		t.Fatalf("unexpected error after retry: %v", err)
	}
	_ = resp.Body.Close()
	if calls != 2 {
		t.Fatalf("expected exactly 2 calls (1 retry), got %d", calls)
	}
}

func TestDoStream_StatusNonRetryable_ParseAPIError(t *testing.T) {
	// Return a non-retryable status (e.g., 418). doStream should call parseAPIError and return *APIError.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
		_, _ = w.Write([]byte(`{"request_id":"rid-parse","error":{"type":"teapot","message":"short and stout"}}`))
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	_, err := c.doStream(context.Background(), http.MethodGet, "/teapot", nil, nil)
	if err == nil {
		t.Fatalf("expected API error")
	}
	if e, ok := IsAPIError(err); !ok || e.RequestID != "rid-parse" || e.Type != "teapot" {
		t.Fatalf("expected APIError with rid-parse/teapot, got %#v", err)
	}
}
