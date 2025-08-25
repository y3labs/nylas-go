package nylas

import (
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
