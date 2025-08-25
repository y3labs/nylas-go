package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestDoJSON_GET_Headers_Query_ResponseHeaders(t *testing.T) {
	// Echo server that checks headers and query; returns a typed envelope.
	type payload struct {
		OK bool `json:"ok"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/foo" {
			t.Fatalf("expected path /foo, got %s", r.URL.Path)
		}
		q := r.URL.Query().Get("query")
		if q != "param" {
			t.Fatalf("expected query=param, got %q", q)
		}

		// Default headers set by Client.do
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("expected Authorization Bearer test-key, got %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected Accept=application/json, got %q", got)
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Fatalf("expected non-empty User-Agent")
		}

		w.Header().Set("X-Test-Header", "test")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-123",
			"data":       payload{OK: true},
		})
	}))
	defer srv.Close()

	c := NewClient("test-key", WithServerURL(srv.URL), WithUserAgent("nylas-go/test"))

	q := url.Values{"query": []string{"param"}}
	out, headers, err := DoJSON[Response[payload]](c, context.Background(), http.MethodGet, "/foo", q, nil, nil)
	if err != nil {
		t.Fatalf("DoJSON error: %v", err)
	}

	if out.RequestID != "rid-123" {
		t.Fatalf("want request_id rid-123, got %s", out.RequestID)
	}
	if out.Data.OK != true {
		t.Fatalf("want OK=true, got %v", out.Data.OK)
	}
	if headers.Get("X-Test-Header") != "test" {
		t.Fatalf("want X-Test-Header=test, got %q", headers.Get("X-Test-Header"))
	}
}

func TestDoJSON_POST_JSONBody(t *testing.T) {
	type payload struct {
		OK bool `json:"ok"`
	}
	type createReq struct {
		Foo string `json:"foo"`
	}

	sawBody := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("expected application/json, got %q", ct)
		}
		var rec createReq
		if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
			t.Fatalf("bad json: %v", err)
		}
		if rec.Foo != "bar" {
			t.Fatalf("expected foo=bar, got %q", rec.Foo)
		}
		sawBody = true
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-456",
			"data":       payload{OK: true},
		})
	}))
	defer srv.Close()

	c := NewClient("test-key", WithServerURL(srv.URL))
	req := createReq{Foo: "bar"}
	out, _, err := DoJSON[Response[payload]](c, context.Background(), http.MethodPost, "/foo", nil, req, nil)
	if err != nil {
		t.Fatalf("DoJSON error: %v", err)
	}
	if !sawBody {
		t.Fatalf("server did not receive JSON body")
	}
	if out.RequestID != "rid-456" || !out.Data.OK {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestDoJSON_ParseAPIError_400(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "123",
			"error": map[string]any{
				"type":           "api_error",
				"message":        "The request is invalid.",
				"provider_error": map[string]any{"foo": "bar"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithServerURL(srv.URL))
	_, _, err := DoJSON[Response[map[string]any]](c, context.Background(), http.MethodGet, "/bad", nil, nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if e, ok := IsAPIError(err); ok {
		if e.Type != "api_error" || e.Message != "The request is invalid." || e.RequestID != "123" || e.StatusCode != 400 {
			t.Fatalf("unexpected api error: %+v", e)
		}
	} else {
		t.Fatalf("want APIError, got %T", err)
	}
}

func TestDoJSON_ParseOAuthError_401(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_request",
			"error_description": "The request is invalid.",
			"error_uri":         "https://docs.nylas.com/reference#authentication-errors",
			"error_code":        100241,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithServerURL(srv.URL))
	_, _, err := DoJSON[Response[map[string]any]](c, context.Background(), http.MethodGet, "/auth", nil, nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if e, ok := IsOAuthError(err); ok {
		if e.ErrorType != "invalid_request" || e.ErrorCode != 100241 || e.ErrorDescription != "The request is invalid." || e.StatusCode != 401 {
			t.Fatalf("unexpected oauth error: %+v", e)
		}
	} else {
		t.Fatalf("want OAuthError, got %T", err)
	}
}

func TestDoJSON_TimeoutWrapsAsSDKTimeoutError(t *testing.T) {
	// slow server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "rid", "data": map[string]any{}})
	}))
	defer srv.Close()

	// very short timeout -> should trigger SDKTimeoutError via wrapTransportError
	hc := &http.Client{Timeout: 10 * time.Millisecond}
	c := NewClient("key", WithServerURL(srv.URL), WithHTTPClient(hc))

	_, _, err := DoJSON[Response[map[string]any]](c, context.Background(), http.MethodGet, "/slow", nil, nil, nil)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if _, ok := IsTimeoutError(err); !ok {
		t.Fatalf("expected SDKTimeoutError, got %T", err)
	}
}
