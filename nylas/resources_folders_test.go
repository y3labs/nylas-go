package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

// --- tests ---

func TestFolderDeserialization(t *testing.T) {
	js := []byte(`{
		"id":"SENT",
		"grant_id":"41009df5-bf11-4c97-aa18-b285b5f2e386",
		"name":"SENT",
		"system_folder":true,
		"object":"folder",
		"unread_count":0,
		"child_count":0,
		"parent_id":"ascsf21412",
		"background_color":"#039BE5",
		"text_color":"#039BE5",
		"total_count":0,
		"attributes":["\\Sent"]
	}`)
	var f models.Folder
	if err := json.Unmarshal(js, &f); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if f.ID != "SENT" ||
		f.GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" ||
		f.Name != "SENT" ||
		f.Object != "folder" {
		t.Fatalf("basic fields mismatch: %#v", f)
	}
	if f.SystemFolder == nil || !*f.SystemFolder {
		t.Fatalf("system_folder = %#v, want true", f.SystemFolder)
	}
	if f.UnreadCount == nil || *f.UnreadCount != 0 {
		t.Fatalf("unread_count = %#v, want 0", f.UnreadCount)
	}
	if f.ChildCount == nil || *f.ChildCount != 0 {
		t.Fatalf("child_count = %#v, want 0", f.ChildCount)
	}
	if f.ParentID == nil || *f.ParentID != "ascsf21412" {
		t.Fatalf("parent_id = %#v, want ascsf21412", f.ParentID)
	}
	if f.BackgroundColor == nil || *f.BackgroundColor != "#039BE5" {
		t.Fatalf("background_color = %#v", f.BackgroundColor)
	}
	if f.TextColor == nil || *f.TextColor != "#039BE5" {
		t.Fatalf("text_color = %#v", f.TextColor)
	}
	if f.TotalCount == nil || *f.TotalCount != 0 {
		t.Fatalf("total_count = %#v, want 0", f.TotalCount)
	}
	if len(f.Attributes) != 1 || f.Attributes[0] != `\Sent` {
		t.Fatalf("attributes = %#v, want [\"\\\\Sent\"]", f.Attributes)
	}
}

func TestListFolders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/folders")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       []any{map[string]any{"id": "folder-1", "grant_id": "abc-123", "name": "Inbox"}},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().List(context.Background(), "abc-123", nil)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListFoldersWithQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/folders")
		q := r.URL.Query()
		if q.Get("limit") != "20" {
			t.Fatalf("limit query mismatch: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc-123", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().List(context.Background(), "abc-123", &ListFoldersParams{Limit: intptr(20)})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListFoldersIncludeHiddenAndSingleLevel(t *testing.T) {
	tests := []struct {
		name string
		p    *ListFoldersParams
		want url.Values
	}{
		{
			"include_hidden=true",
			&ListFoldersParams{IncludeHiddenFolders: boolptr(true)},
			url.Values{"include_hidden_folders": []string{"true"}},
		},
		{
			"single_level=true",
			&ListFoldersParams{SingleLevel: boolptr(true)},
			url.Values{"single_level": []string{"true"}},
		},
		{
			"include_hidden=false",
			&ListFoldersParams{IncludeHiddenFolders: boolptr(false)},
			url.Values{"include_hidden_folders": []string{"false"}},
		},
		{
			"single_level=false",
			&ListFoldersParams{SingleLevel: boolptr(false)},
			url.Values{"single_level": []string{"false"}},
		},
		{
			"multiple params",
			&ListFoldersParams{
				Limit:                intptr(20),
				ParentID:             strptr("parent-123"),
				IncludeHiddenFolders: boolptr(true),
				SingleLevel:          boolptr(true),
			},
			url.Values{
				"limit":                  []string{"20"},
				"parent_id":              []string{"parent-123"},
				"include_hidden_folders": []string{"true"},
				"single_level":           []string{"true"},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/folders")
				q := r.URL.Query()
				for k, vs := range tc.want {
					if q.Get(k) != vs[0] {
						t.Fatalf("want %s=%q, got %q (raw=%s)", k, vs[0], q.Get(k), r.URL.RawQuery)
					}
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc-123", "data": []any{}})
			}))
			defer ts.Close()

			c := newTestClient(ts.URL, "test-key")
			if _, err := c.Folders().List(context.Background(), "abc-123", tc.p); err != nil {
				t.Fatalf("List error: %v", err)
			}
		})
	}
}

func TestListFoldersWithSelectParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/folders")
		if r.URL.Query().Get("select") != "id,name,total_count,unread_count" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []any{map[string]any{
				"id":           "folder-123",
				"name":         "Important",
				"total_count":  42,
				"unread_count": 5,
			}},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().List(context.Background(), "abc-123", &ListFoldersParams{
		Select: strptr("id,name,total_count,unread_count"),
	})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindFolder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/folders/folder-123")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "folder-123", "grant_id": "abc-123", "name": "Important"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().Get(context.Background(), "abc-123", "folder-123", nil)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestFindFolderWithSelectParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/folders/folder-123")
		if r.URL.Query().Get("select") != "id,name,total_count,unread_count" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "folder-123", "name": "Important", "total_count": 42, "unread_count": 5},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().Get(context.Background(), "abc-123", "folder-123", &FindFolderParams{
		Select: strptr("id,name,total_count,unread_count"),
	})
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestCreateFolder(t *testing.T) {
	wantBody := map[string]any{
		"name":             "My New Folder",
		"parent_id":        "parent-folder-id",
		"background_color": "#039BE5",
		"text_color":       "#039BE5",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/folders")
		var got map[string]any
		_ = json.NewDecoder(r.Body).Decode(&got)
		for k, v := range wantBody {
			if got[k] != v {
				t.Fatalf("body[%s] = %#v, want %#v", k, got[k], v)
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "folder-123", "name": "My New Folder", "grant_id": "abc-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().Create(context.Background(), "abc-123", CreateFolderRequest{
		Name:            "My New Folder",
		ParentID:        strptr("parent-folder-id"),
		BackgroundColor: strptr("#039BE5"),
		TextColor:       strptr("#039BE5"),
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
}

func TestUpdateFolder(t *testing.T) {
	wantBody := map[string]any{
		"name":             "My New Folder",
		"parent_id":        "parent-folder-id",
		"background_color": "#039BE5",
		"text_color":       "#039BE5",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// UPDATE uses PATCH in Go implementation
		assertMethodPath(t, r, http.MethodPatch, "/v3/grants/abc-123/folders/folder-123")
		var got map[string]any
		_ = json.NewDecoder(r.Body).Decode(&got)
		for k, v := range wantBody {
			if got[k] != v {
				t.Fatalf("body[%s] = %#v, want %#v", k, got[k], v)
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "folder-123", "name": "My New Folder", "grant_id": "abc-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Folders().Update(context.Background(), "abc-123", "folder-123", UpdateFolderRequest{
		Name:            strptr("My New Folder"),
		ParentID:        strptr("parent-folder-id"),
		BackgroundColor: strptr("#039BE5"),
		TextColor:       strptr("#039BE5"),
	})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestDestroyFolder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/folders/folder-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if err := c.Folders().Delete(context.Background(), "abc-123", "folder-123"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}
