package nylas

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseAPIError_OAuth(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Code = 400
	rec.Header().Set("x-request-id", "rid-o")
	rec.Body.WriteString(`{"error":"invalid_grant","error_code":123,"error_description":"bad","error_uri":"https://x"}`)

	err := parseAPIError(rec.Result())
	e, ok := IsOAuthError(err)
	if !ok || e.ErrorType != "invalid_grant" || e.ErrorCode != 123 || e.RequestID != "rid-o" || !IsStatus(err, 400) {
		t.Fatalf("bad oauth error: %#v", err)
	}
}

func TestParseAPIError_Envelope(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Code = 429
	rec.Header().Set("x-request-id", "rid-a")
	rec.Body.WriteString(`{"request_id":"abc","error":{"type":"rate_limit","message":"slow","provider_error":"foo"}}`)

	err := parseAPIError(rec.Result())
	e, ok := IsAPIError(err)
	if !ok || e.Type != "rate_limit" || e.Message != "slow" || e.RequestID != "abc" || e.ProviderErrorString() != "foo" || !IsStatus(err, 429) {
		t.Fatalf("bad api error: %#v", err)
	}
}

func TestParseAPIError_FlatAndFallback(t *testing.T) {
	// Flat
	{
		rec := httptest.NewRecorder()
		rec.Code = 403
		rec.Header().Set("x-request-id", "rid-flat")
		rec.Body.WriteString(`{"message":"nope","type":"forbidden"}`)
		err := parseAPIError(rec.Result())
		e, ok := IsAPIError(err)
		if !ok || e.RequestID != "rid-flat" || e.Message != "nope" || e.Type != "forbidden" || !IsStatus(err, 403) {
			t.Fatalf("flat: %#v", err)
		}
	}
	// Fallback (non-JSON)
	{
		rec := httptest.NewRecorder()
		rec.Code = 500
		rec.Body.WriteString(`internal server error`)
		err := parseAPIError(rec.Result())
		e, ok := IsAPIError(err)
		if !ok || e.StatusCode != 500 || !strings.Contains(e.Message, "Internal Server Error") {
			t.Fatalf("fallback: %#v", err)
		}
	}
}

func TestProviderErrorString_ObjectOrString(t *testing.T) {
	// string
	{
		e := &APIError{ProviderError: []byte(`"oops"`)}
		if e.ProviderErrorString() != "oops" {
			t.Fatalf("want string")
		}
	}
	// object
	{
		e := &APIError{ProviderError: []byte(`{"msg":"oops"}`)}
		if e.ProviderErrorString() != `{"msg":"oops"}` {
			t.Fatalf("want raw object")
		}
	}
}

func TestWrapTransportError(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://x", nil)

	// context deadline exceeded
	err := wrapTransportError(context.DeadlineExceeded, req, 2*time.Second)
	if _, ok := IsTimeoutError(err); !ok {
		t.Fatalf("want SDKTimeoutError for context deadline")
	}

	// net.Error timeout
	tErr := &timeoutErr{}
	err = wrapTransportError(tErr, req, time.Second)
	if _, ok := IsTimeoutError(err); !ok {
		t.Fatalf("want SDKTimeoutError for net.Error timeout")
	}

	// passthrough
	raw := errors.New("boom")
	if got := wrapTransportError(raw, req, 0); !errors.Is(got, raw) {
		t.Fatalf("want passthrough")
	}
}

func TestErrorStrings(t *testing.T) {
	ae := &APIError{Message: "m", Type: "t", StatusCode: 418}
	if s := ae.Error(); !strings.Contains(s, "nylas: m") || !strings.Contains(s, "status=418") {
		t.Fatalf("APIError.Error: %q", s)
	}
	ae.RequestID = "rid"
	if s := ae.Error(); !strings.Contains(s, "rid=rid") {
		t.Fatalf("APIError.Error with rid: %q", s)
	}

	oe := &OAuthError{ErrorDescription: "bad", ErrorCode: 123, StatusCode: 400}
	if s := oe.Error(); !strings.Contains(s, "nylas-oauth: bad") || !strings.Contains(s, "code=123") {
		t.Fatalf("OAuthError.Error: %q", s)
	}
	oe.RequestID = "rid2"
	if s := oe.Error(); !strings.Contains(s, "rid=rid2") {
		t.Fatalf("OAuthError.Error with rid: %q", s)
	}

	se := &SDKTimeoutError{URL: "http://x", Timeout: 2 * time.Second}
	if s := se.Error(); !strings.Contains(s, "timed out") || !strings.Contains(s, "http://x") {
		t.Fatalf("SDKTimeoutError.Error: %q", s)
	}
}

func TestIsStatus(t *testing.T) {
	ae := &APIError{StatusCode: http.StatusNotFound}
	if !IsStatus(ae, http.StatusNotFound) || IsStatus(ae, http.StatusOK) {
		t.Fatalf("IsStatus failed for APIError")
	}
	oe := &OAuthError{StatusCode: http.StatusUnauthorized}
	if !IsStatus(oe, http.StatusUnauthorized) || IsStatus(oe, http.StatusOK) {
		t.Fatalf("IsStatus failed for OAuthError")
	}
}
