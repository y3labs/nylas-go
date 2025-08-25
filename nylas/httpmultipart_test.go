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
