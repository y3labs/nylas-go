package nylas

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func strptr(s string) *string { return &s }
func ptr[T any](v T) *T       { return &v }
func boolptr(b bool) *bool    { return &b }
func i64ptr(v int64) *int64   { return &v }
func intptr(i int) *int       { return &i }

func strval(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}

func boolval(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func assertStrPtr(t *testing.T, got *string, want, field string) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("%s = %q, want %q (nil? %v)", field, strval(got), want, got == nil)
	}
}

func assertMethodPath(t *testing.T, r *http.Request, method, wantPath string) {
	t.Helper()
	if r.Method != method {
		t.Fatalf("method = %s, want %s", r.Method, method)
	}
	gotEsc := r.URL.EscapedPath()
	if gotEsc == wantPath {
		return
	}
	// fallback: compare decoded forms
	wantDec, _ := url.PathUnescape(wantPath)
	if r.URL.Path != wantDec {
		t.Fatalf("path = %s (escaped %s), want %s", r.URL.Path, gotEsc, wantPath)
	}
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) { return 0, errors.New("simulated copy error") }

// rtErr is a RoundTripper that always returns the provided error.
type rtErr struct{ err error }

func (r rtErr) RoundTrip(*http.Request) (*http.Response, error) { return nil, r.err }

// timeoutErr implements net.Error with Timeout() == true (so wrapTransportError produces SDKTimeoutError).
type timeoutErr struct{}

func (t timeoutErr) Error() string   { return "simulated timeout" }
func (t timeoutErr) Timeout() bool   { return true }
func (t timeoutErr) Temporary() bool { return false }

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (len(sub) == 0 || (stringIndex(s, sub) >= 0))
}
func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestPtr(t *testing.T) {
	val := 42
	ptr := Ptr(val)
	if ptr == nil || *ptr != 42 {
		t.Fatalf("expected *ptr=42, got %v", ptr)
	}
}

func TestStringPtr(t *testing.T) {
	s := "hello"
	ptr := StringPtr(s)
	if ptr == nil || *ptr != "hello" {
		t.Fatalf("expected *ptr=hello, got %v", ptr)
	}
}
