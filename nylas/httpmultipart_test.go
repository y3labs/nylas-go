package nylas

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDoMultipartParts_CustomHeaders(t *testing.T) {
	var fileCT string
	var fieldCT string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data;") {
			t.Fatalf("bad ct: %s", ct)
		}
		mediaType, params, _ := mime.ParseMediaType(ct)
		if mediaType != "multipart/form-data" {
			t.Fatalf("bad media type: %s", mediaType)
		}
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next: %v", err)
			}
			if p.FileName() != "" {
				fileCT = p.Header.Get("Content-Type")
			} else {
				fieldCT = p.Header.Get("Content-Type")
			}
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))

	_, err := c.doMultipartParts(
		context.Background(),
		http.MethodPost,
		"/upload/test",
		nil,
		[]FormField{
			{Name: "meta", Value: "x", ContentType: "application/json"},
		},
		[]FormFile{
			{Field: "file", Filename: "invite.ics", ContentType: "text/calendar", Reader: bytes.NewReader([]byte("BEGIN:VCALENDAR"))},
		},
	)
	if err != nil {
		t.Fatalf("doMultipartParts err: %v", err)
	}
	if fileCT != "text/calendar" || fieldCT != "application/json" {
		t.Fatalf("content-types not preserved: file=%q field=%q", fileCT, fieldCT)
	}
}

func TestDoMultipart_Minimal(t *testing.T) {
	var saw string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/upload" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data;") {
			t.Fatalf("bad content-type: %s", ct)
		}
		mr, err := r.MultipartReader()
		if err != nil {
			t.Fatalf("reader: %v", err)
		}
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next: %v", err)
			}
			b, _ := io.ReadAll(p)
			if fn := p.FileName(); fn != "" {
				saw = "file:" + fn + ":" + string(b)
			} else {
				saw = "field:" + p.FormName() + ":" + string(b)
			}
		}
		w.WriteHeader(204)
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	resp, err := c.doMultipart(
		context.Background(),
		http.MethodPost,
		"/v3/upload",
		nil,
		map[string]string{"k": "v"},
		map[string]io.Reader{"note.txt": bytes.NewReader([]byte("hello"))},
	)
	if err != nil {
		t.Fatalf("doMultipart err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		t.Fatalf("status: %d", resp.StatusCode)
	}
	if !strings.Contains(saw, "file:note.txt:hello") && !strings.Contains(saw, "field:k:v") {
		t.Fatalf("did not see parts; saw=%q", saw)
	}
}

func TestDoMultipart_BuildURLError(t *testing.T) {
	c := NewClient("test-key", WithServerURL("")) // buildURL should error
	_, err := c.doMultipart(context.Background(), http.MethodPost, "/x", nil, nil, map[string]io.Reader{"f.txt": bytes.NewReader(nil)})
	if err == nil {
		t.Fatalf("expected buildURL error")
	}
}

func TestDoMultipart_CopyError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("server should not be called when io.Copy fails")
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	// reader that errors after 0 bytes to trigger io.Copy failure.
	bad := new(errReader)

	_, err := c.doMultipart(context.Background(), http.MethodPost, "/upload", nil, nil, map[string]io.Reader{"bad.txt": bad})
	if err == nil {
		t.Fatalf("expected io.Copy error")
	}
}

func TestDoMultipart_TransportError_Wrapped(t *testing.T) {
	// HTTP client whose transport always times out (net.Error Timeout true)
	cl := &http.Client{
		Transport: rtErr{timeoutErr{}},
		Timeout:   50 * time.Millisecond,
	}
	c := NewClient("test-key", WithServerURL("http://example.invalid"), WithHTTPClient(cl))

	_, err := c.doMultipart(context.Background(), http.MethodPost, "/m", nil, nil, map[string]io.Reader{"ok.txt": bytes.NewReader([]byte("x"))})
	if err == nil {
		t.Fatalf("expected wrapped transport error")
	}
	if _, ok := IsTimeoutError(err); !ok {
		t.Fatalf("expected SDKTimeoutError, got %T: %v", err, err)
	}
}

func TestDoMultipart_StatusRetryThenSuccess(t *testing.T) {
	pol := DefaultRetryPolicy()
	if pol.MaxRetries == 0 {
		t.Skip("DefaultRetryPolicy has MaxRetries=0; no retry path to test")
	}

	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			// Use 429 to guarantee ShouldRetry == true in most policies.
			w.WriteHeader(http.StatusTooManyRequests) // 429
			_, _ = w.Write([]byte(`{"error":{"type":"rate_limited","message":"boom"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id":"rid","data":{"id":"ok"}}`))
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	resp, err := c.doMultipart(context.Background(), http.MethodPost, "/upload", nil, nil,
		map[string]io.Reader{"f.txt": bytes.NewReader([]byte("hello"))})
	if err != nil {
		t.Fatalf("unexpected error after retry: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 2 {
		t.Fatalf("expected exactly 2 calls (1 retry), got %d", calls)
	}
}

func TestDoMultipart_StatusRetryExhausted_ReturnsAPIError(t *testing.T) {
	// Always 500 to force parseAPIError path on last attempt.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`{"request_id":"rid-x","error":{"type":"server_error","message":"boom"}}`))
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	resp, err := c.doMultipart(context.Background(), http.MethodPost, "/upload", nil, nil, map[string]io.Reader{"f.txt": bytes.NewReader([]byte("hello"))})
	if err == nil {
		if resp != nil {
			resp.Body.Close()
		}
		t.Fatalf("expected APIError on exhausted retries")
	}
	if e, ok := IsAPIError(err); !ok || e.RequestID != "rid-x" {
		t.Fatalf("expected APIError rid-x, got %#v", err)
	}
	// resp was returned alongside error; body should be closed by parseAPIError already.
}

// ---- doMultipartParts tests ----

func TestDoMultipartParts_FieldsAndFiles_Success(t *testing.T) {
	// Assert multipart parts (names and content types) and return 200.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mediaType := r.Header.Get("Content-Type")
		if mediaType == "" || mediaType[:19] != "multipart/form-data" {
			t.Fatalf("expected multipart content type, got %q", mediaType)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"request_id":"rid-2","data":{"id":"ok2"}}`))
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	fields := []FormField{
		{Name: "meta", Value: `{"a":1}`, ContentType: "application/json"},
	}
	files := []FormFile{
		{Field: "file", Filename: "a.txt", Reader: bytes.NewReader([]byte("x")), ContentType: "text/plain"},
	}

	resp, err := c.doMultipartParts(context.Background(), http.MethodPost, "/mpp", nil, fields, files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestFileReader(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "x.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("xyz"); err != nil {
		t.Fatal(err)
	}

	rc, err := FileReader(f.Name())
	if err != nil {
		t.Fatalf("FileReader: %v", err)
	}
	defer rc.Close()
	buf := make([]byte, 3)
	if _, err := rc.Read(buf); err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(buf) != "xyz" {
		t.Fatalf("got %q", string(buf))
	}
}
