package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// test item used as a resource model
type testItem struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// mockResource mirrors how real resources are implemented (calls DoJSON under the hood).
type mockResource struct{ c *Client }

func (r *mockResource) List(ctx context.Context) (*ListResponse[testItem], error) {
	out, headers, err := DoJSON[ListResponse[testItem]](r.c, ctx, http.MethodGet, "/foo", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

func (r *mockResource) Find(ctx context.Context, id string) (*Response[testItem], error) {
	path := "/foo/" + url.PathEscape(id)
	out, headers, err := DoJSON[Response[testItem]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

func (r *mockResource) Create(ctx context.Context, body testItem) (*Response[testItem], error) {
	out, headers, err := DoJSON[Response[testItem]](r.c, ctx, http.MethodPost, "/foo", nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

func (r *mockResource) Update(ctx context.Context, id string, body map[string]any) (*Response[testItem], error) {
	path := "/foo/" + url.PathEscape(id)
	out, headers, err := DoJSON[Response[testItem]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

func (r *mockResource) Destroy(ctx context.Context, id string) (*DeleteResponse, error) {
	path := "/foo/" + url.PathEscape(id)
	out, headers, err := DoJSON[DeleteResponse](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

func TestResource_List_Find_Create_Update_Destroy_WithHeaders(t *testing.T) {
	var lastMethod, lastPath string
	var sawCreateBody, sawUpdateBody bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastMethod, lastPath = r.Method, r.URL.Path
		// common header in all responses
		w.Header().Set("X-Test-Header", "test")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/foo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id": "rid-list",
				"data":       []testItem{{ID: "1", Name: "a"}, {ID: "2", Name: "b"}},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/foo/123":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id": "rid-find",
				"data":       testItem{ID: "123", Name: "n"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/foo":
			var body testItem
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body.Name != "created" {
				t.Fatalf("expected create body name=created, got %+v", body)
			}
			sawCreateBody = true
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id": "rid-create",
				"data":       testItem{ID: "999", Name: body.Name},
			})
		case r.Method == http.MethodPut && r.URL.Path == "/foo/777":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["name"] != "updated" {
				t.Fatalf("expected update body name=updated, got %+v", body)
			}
			sawUpdateBody = true
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id": "rid-update",
				"data":       testItem{ID: "777", Name: "updated"},
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/foo/777":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id": "rid-delete",
				"data":       map[string]any{"status": "ok"},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := NewClient("key", WithServerURL(srv.URL))
	r := &mockResource{c: c}
	ctx := context.Background()

	// List
	lr, err := r.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if lastMethod != http.MethodGet || lastPath != "/foo" {
		t.Fatalf("list wrong call: %s %s", lastMethod, lastPath)
	}
	if lr.RequestID != "rid-list" || len(lr.Data) != 2 {
		t.Fatalf("bad list response: %+v", lr)
	}
	if lr.Headers.Get("X-Test-Header") != "test" {
		t.Fatalf("list missing header propagation")
	}

	// Find
	fr, err := r.Find(ctx, "123")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if lastMethod != http.MethodGet || lastPath != "/foo/123" {
		t.Fatalf("find wrong call: %s %s", lastMethod, lastPath)
	}
	if fr.RequestID != "rid-find" || fr.Data.ID != "123" {
		t.Fatalf("bad find response: %+v", fr)
	}
	if fr.Headers.Get("X-Test-Header") != "test" {
		t.Fatalf("find missing header propagation")
	}

	// Create
	cr, err := r.Create(ctx, testItem{Name: "created"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !sawCreateBody {
		t.Fatalf("create body not seen by server")
	}
	if lastMethod != http.MethodPost || lastPath != "/foo" {
		t.Fatalf("create wrong call: %s %s", lastMethod, lastPath)
	}
	if cr.RequestID != "rid-create" || cr.Data.ID != "999" {
		t.Fatalf("bad create response: %+v", cr)
	}
	if cr.Headers.Get("X-Test-Header") != "test" {
		t.Fatalf("create missing header propagation")
	}

	// Update
	ur, err := r.Update(ctx, "777", map[string]any{"name": "updated"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !sawUpdateBody {
		t.Fatalf("update body not seen by server")
	}
	if lastMethod != http.MethodPut || lastPath != "/foo/777" {
		t.Fatalf("update wrong call: %s %s", lastMethod, lastPath)
	}
	if ur.RequestID != "rid-update" || ur.Data.ID != "777" {
		t.Fatalf("bad update response: %+v", ur)
	}
	if ur.Headers.Get("X-Test-Header") != "test" {
		t.Fatalf("update missing header propagation")
	}

	// Destroy
	dr, err := r.Destroy(context.Background(), "777")
	if err != nil {
		t.Fatalf("destroy: %v", err)
	}
	if lastMethod != http.MethodDelete || lastPath != "/foo/777" {
		t.Fatalf("destroy wrong call: %s %s", lastMethod, lastPath)
	}
	if dr.RequestID != "rid-delete" {
		t.Fatalf("bad delete response: %+v", dr)
	}
	if dr.Headers.Get("X-Test-Header") != "test" {
		t.Fatalf("destroy missing header propagation")
	}
}
