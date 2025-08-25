package nylas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// -----------------------------
// Types
// -----------------------------

// APIError models the general Nylas error envelope:
// {"request_id":"...", "error":{"type":"...","message":"...","provider_error":{...}}}
type APIError struct {
	RequestID     string          `json:"request_id"`
	StatusCode    int             `json:"-"`
	Type          string          `json:"type,omitempty"`
	Message       string          `json:"message,omitempty"`
	ProviderError json.RawMessage `json:"provider_error,omitempty"` // may be string OR object
	RawBody       string          `json:"-"`
	Headers       http.Header     `json:"-"`
	RateLimit     RateLimitInfo   `json:"-"`
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("nylas: %s (type=%s status=%d rid=%s)", e.Message, e.Type, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("nylas: %s (type=%s status=%d)", e.Message, e.Type, e.StatusCode)
}

// ProviderErrorString is a convenience to display provider_error even if it was an object.
func (e *APIError) ProviderErrorString() string {
	if len(e.ProviderError) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(e.ProviderError, &s); err == nil {
		return s
	}
	return string(e.ProviderError)
}

// OAuthError models OAuth error responses:
// {"error":"...","error_code":123,"error_description":"...","error_uri":"..."}
type OAuthError struct {
	RequestID        string        `json:"request_id,omitempty"`
	StatusCode       int           `json:"-"`
	ErrorType        string        `json:"error"`
	ErrorCode        int           `json:"error_code"`
	ErrorDescription string        `json:"error_description"`
	ErrorURI         string        `json:"error_uri"`
	RawBody          string        `json:"-"`
	Headers          http.Header   `json:"-"`
	RateLimit        RateLimitInfo `json:"-"`
}

func (e *OAuthError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("nylas-oauth: %s (code=%d status=%d rid=%s)", e.ErrorDescription, e.ErrorCode, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("nylas-oauth: %s (code=%d status=%d)", e.ErrorDescription, e.ErrorCode, e.StatusCode)
}

// SDKTimeoutError models client-side timeouts (matches Python NylasSdkTimeoutError intent).
type SDKTimeoutError struct {
	URL     string
	Timeout time.Duration
	Headers http.Header
}

func (e *SDKTimeoutError) Error() string {
	return fmt.Sprintf("nylas SDK timed out before receiving a response (url=%s timeout=%s)", e.URL, e.Timeout)
}

// -----------------------------
// Helpers to classify errors
// -----------------------------

func IsAPIError(err error) (*APIError, bool) {
	var e *APIError
	ok := errors.As(err, &e)
	return e, ok
}

func IsOAuthError(err error) (*OAuthError, bool) {
	var e *OAuthError
	ok := errors.As(err, &e)
	return e, ok
}

func IsTimeoutError(err error) (*SDKTimeoutError, bool) {
	var e *SDKTimeoutError
	ok := errors.As(err, &e)
	return e, ok
}

// Common status helpers (optional sugar)
func IsStatus(err error, code int) bool {
	if e, ok := IsAPIError(err); ok && e.StatusCode == code {
		return true
	}
	if e, ok := IsOAuthError(err); ok && e.StatusCode == code {
		return true
	}
	return false
}

// -----------------------------
// Parsing non-2xx HTTP responses
// -----------------------------

// parseAPIError reads an error response body and returns an *APIError or *OAuthError.
// Always returns a typed error (never nil).
func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	// Try OAuth shape first
	var oauth struct {
		RequestID        string `json:"request_id"`
		Error            string `json:"error"`
		ErrorCode        int    `json:"error_code"`
		ErrorDescription string `json:"error_description"`
		ErrorURI         string `json:"error_uri"`
	}
	if json.Unmarshal(b, &oauth) == nil && oauth.Error != "" && oauth.ErrorDescription != "" {
		return &OAuthError{
			RequestID:        firstNonEmpty(oauth.RequestID, resp.Header.Get("x-request-id")),
			StatusCode:       resp.StatusCode,
			ErrorType:        oauth.Error,
			ErrorCode:        oauth.ErrorCode,
			ErrorDescription: oauth.ErrorDescription,
			ErrorURI:         oauth.ErrorURI,
			RawBody:          string(b),
			Headers:          resp.Header,
			RateLimit:        ParseRateLimit(resp.Header),
		}
	}

	// Try general envelope
	var envelope struct {
		RequestID string `json:"request_id"`
		Error     struct {
			Type          string          `json:"type"`
			Message       string          `json:"message"`
			ProviderError json.RawMessage `json:"provider_error"`
		} `json:"error"`
	}
	if json.Unmarshal(b, &envelope) == nil && (envelope.RequestID != "" || envelope.Error.Message != "" || envelope.Error.Type != "") {
		return &APIError{
			RequestID:     firstNonEmpty(envelope.RequestID, resp.Header.Get("x-request-id")),
			StatusCode:    resp.StatusCode,
			Type:          envelope.Error.Type,
			Message:       envelope.Error.Message,
			ProviderError: envelope.Error.ProviderError,
			RawBody:       string(b),
			Headers:       resp.Header,
			RateLimit:     ParseRateLimit(resp.Header),
		}
	}

	// As a fallback, try to bind directly to APIError fields at top-level (some services respond that way)
	var flat APIError
	if json.Unmarshal(b, &flat) == nil && (flat.Message != "" || flat.Type != "" || flat.RequestID != "") {
		flat.StatusCode = resp.StatusCode
		flat.RawBody = string(b)
		flat.Headers = resp.Header
		flat.RateLimit = ParseRateLimit(resp.Header)
		if flat.RequestID == "" {
			flat.RequestID = resp.Header.Get("x-request-id")
		}
		return &flat
	}

	// Last resort: synthesize minimal APIError
	return &APIError{
		RequestID:  resp.Header.Get("x-request-id"),
		StatusCode: resp.StatusCode,
		Message:    http.StatusText(resp.StatusCode),
		RawBody:    string(b),
		Headers:    resp.Header,
		RateLimit:  ParseRateLimit(resp.Header),
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// -----------------------------
// Client-side wrapping (timeouts)
// -----------------------------

// wrapTransportError converts context/transport timeouts into SDKTimeoutError.
// Call this in your Client HTTP path where you return errors from Do/DoJSON.
func wrapTransportError(err error, req *http.Request, configuredTimeout time.Duration) error {
	if err == nil {
		return nil
	}
	// context deadline exceeded
	if errors.Is(err, context.DeadlineExceeded) {
		return &SDKTimeoutError{URL: req.URL.String(), Timeout: configuredTimeout}
	}
	// net.Error with Timeout
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return &SDKTimeoutError{URL: req.URL.String(), Timeout: configuredTimeout}
	}
	return err
}
