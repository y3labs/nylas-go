package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestContactDeserialization(t *testing.T) {
	js := []byte(`{
		"birthday": "1960-12-31",
		"company_name": "Nylas",
		"emails": [{"type": "work", "email": "john-work@example.com"}],
		"given_name": "John",
		"grant_id": "41009df5-bf11-4c97-aa18-b285b5f2e386",
		"groups": [{"id": "starred"}],
		"id": "5d3qmne77v32r8l4phyuksl2x",
		"im_addresses": [{"type": "other", "im_address": "myjabberaddress"}],
		"job_title": "Software Engineer",
		"manager_name": "Bill",
		"middle_name": "Jacob",
		"nickname": "JD",
		"notes": "Loves ramen",
		"object": "contact",
		"office_location": "123 Main Street",
		"phone_numbers": [{"type": "work", "number": "+1-555-555-5555"}],
		"physical_addresses": [{
			"type": "work",
			"street_address": "123 Main Street",
			"postal_code": "94107",
			"state": "CA",
			"country": "US",
			"city": "San Francisco"
		}],
		"picture_url": "https://example.com/picture.jpg",
		"suffix": "Jr.",
		"surname": "Doe",
		"web_pages": [{"type": "work", "url": "http://www.linkedin.com/in/johndoe"}]
	}`)

	var c models.Contact
	if err := json.Unmarshal(js, &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assertStrPtr(t, c.Birthday, "1960-12-31", "Birthday")
	assertStrPtr(t, c.CompanyName, "Nylas", "CompanyName")
	if len(c.Emails) != 1 {
		t.Fatalf("Emails length = %d, want 1", len(c.Emails))
	}
	assertStrPtr(t, c.Emails[0].Email, "john-work@example.com", "Emails[0].Email")
	assertStrPtr(t, c.Emails[0].Type, "work", "Emails[0].Type")

	assertStrPtr(t, c.GivenName, "John", "GivenName")
	if c.GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" {
		t.Fatalf("GrantID = %q", c.GrantID)
	}
	if len(c.Groups) != 1 || c.Groups[0].ID != "starred" {
		t.Fatalf("Groups = %#v", c.Groups)
	}
	if c.ID != "5d3qmne77v32r8l4phyuksl2x" {
		t.Fatalf("ID = %q", c.ID)
	}
	if len(c.IMAddresses) != 1 {
		t.Fatalf("IMAddresses length = %d, want 1", len(c.IMAddresses))
	}
	assertStrPtr(t, c.IMAddresses[0].Type, "other", "IMAddresses[0].Type")
	assertStrPtr(t, c.IMAddresses[0].IMAddress, "myjabberaddress", "IMAddresses[0].IMAddress")

	assertStrPtr(t, c.JobTitle, "Software Engineer", "JobTitle")
	assertStrPtr(t, c.ManagerName, "Bill", "ManagerName")
	assertStrPtr(t, c.MiddleName, "Jacob", "MiddleName")
	assertStrPtr(t, c.Nickname, "JD", "Nickname")
	assertStrPtr(t, c.Notes, "Loves ramen", "Notes")
	if c.Object != "contact" {
		t.Fatalf("Object = %q, want %q", c.Object, "contact")
	}
	assertStrPtr(t, c.OfficeLocation, "123 Main Street", "OfficeLocation")

	if len(c.PhoneNumbers) != 1 {
		t.Fatalf("PhoneNumbers length = %d, want 1", len(c.PhoneNumbers))
	}
	assertStrPtr(t, c.PhoneNumbers[0].Type, "work", "PhoneNumbers[0].Type")
	assertStrPtr(t, c.PhoneNumbers[0].Number, "+1-555-555-5555", "PhoneNumbers[0].Number")

	if len(c.PhysicalAddresses) != 1 {
		t.Fatalf("PhysicalAddresses length = %d, want 1", len(c.PhysicalAddresses))
	}
	addr := c.PhysicalAddresses[0]
	assertStrPtr(t, addr.Type, "work", "PhysicalAddresses[0].Type")
	assertStrPtr(t, addr.StreetAddress, "123 Main Street", "PhysicalAddresses[0].StreetAddress")
	assertStrPtr(t, addr.PostalCode, "94107", "PhysicalAddresses[0].PostalCode")
	assertStrPtr(t, addr.State, "CA", "PhysicalAddresses[0].State")
	assertStrPtr(t, addr.Country, "US", "PhysicalAddresses[0].Country")
	assertStrPtr(t, addr.City, "San Francisco", "PhysicalAddresses[0].City")

	assertStrPtr(t, c.PictureURL, "https://example.com/picture.jpg", "PictureURL")
	assertStrPtr(t, c.Suffix, "Jr.", "Suffix")
	assertStrPtr(t, c.Surname, "Doe", "Surname")

	if len(c.WebPages) != 1 {
		t.Fatalf("WebPages length = %d, want 1", len(c.WebPages))
	}
	assertStrPtr(t, c.WebPages[0].Type, "work", "WebPages[0].Type")
	assertStrPtr(t, c.WebPages[0].URL, "http://www.linkedin.com/in/johndoe", "WebPages[0].URL")
}

func TestListContacts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Contacts().List(context.Background(), "abc-123", nil)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if out == nil {
		t.Fatalf("List returned nil")
	}
}

func TestListContactsWithQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts")
		if g := r.URL.Query().Get("limit"); g != "20" {
			t.Fatalf("limit = %q, want 20", g)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	lim := 20
	q := &models.ListContactsQueryParams{Limit: &lim}
	_, err := c.Contacts().List(context.Background(), "abc-123", q)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListContactsWithSelectParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts")
		if g := r.URL.Query().Get("select"); g != "id,given_name,surname,emails" {
			t.Fatalf("select = %q, want %q", g, "id,given_name,surname,emails")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{"id": "contact-123", "given_name": "John", "surname": "Doe", "emails": []map[string]any{{"email": "john@example.com", "type": "work"}}},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	q := &models.ListContactsQueryParams{Select: strptr("id,given_name,surname,emails")}
	_, err := c.Contacts().List(context.Background(), "abc-123", q)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindContact(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts/contact-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":         "contact-123",
				"grant_id":   "abc-123",
				"given_name": "John",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Contacts().Get(context.Background(), "abc-123", "contact-123", nil)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if out == nil || out.Data.ID != "contact-123" {
		t.Fatalf("unexpected get response: %#v", out)
	}
}

func TestFindContactWithSelectParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts/contact-123")
		if g := r.URL.Query().Get("select"); g != "id,given_name,surname,emails" {
			t.Fatalf("select = %q, want %q", g, "id,given_name,surname,emails")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":         "contact-123",
				"given_name": "John",
				"surname":    "Doe",
				"emails":     []map[string]any{{"email": "john@example.com", "type": "work"}},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	q := &models.FindContactQueryParams{Select: strptr("id,given_name,surname,emails")}
	_, err := c.Contacts().Get(context.Background(), "abc-123", "contact-123", q)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestFindContactWithQueryParams_ProfilePicture(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts/contact-123")
		if g := r.URL.Query().Get("profile_picture"); g != "true" {
			t.Fatalf("profile_picture = %q, want true", g)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "contact-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	pp := true
	q := &models.FindContactQueryParams{ProfilePicture: &pp}
	_, err := c.Contacts().Get(context.Background(), "abc-123", "contact-123", q)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestCreateContact(t *testing.T) {
	var captured map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/contacts")
		captured = decodeJSONBody(t, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":           "contact-123",
				"grant_id":     "abc-123",
				"given_name":   captured["given_name"],
				"surname":      captured["surname"],
				"company_name": captured["company_name"],
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	body := models.CreateContactRequest{
		GivenName:   strptr("John"),
		Surname:     strptr("Doe"),
		CompanyName: strptr("Nylas"),
	}
	_, err := c.Contacts().Create(context.Background(), "abc-123", body)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if captured["given_name"] != "John" || captured["surname"] != "Doe" || captured["company_name"] != "Nylas" {
		t.Fatalf("unexpected posted body: %#v", captured)
	}
}

func TestUpdateContact(t *testing.T) {
	var captured map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/contacts/contact-123")
		captured = decodeJSONBody(t, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "contact-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	body := models.UpdateContactRequest{
		GivenName:   strptr("John"),
		Surname:     strptr("Doe"),
		CompanyName: strptr("Nylas"),
	}
	_, err := c.Contacts().Update(context.Background(), "abc-123", "contact-123", body)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if captured["given_name"] != "John" || captured["surname"] != "Doe" || captured["company_name"] != "Nylas" {
		t.Fatalf("unexpected posted body: %#v", captured)
	}
}

func TestDestroyContact(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/contacts/contact-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc-123"})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Contacts().Delete(context.Background(), "abc-123", "contact-123")
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if out == nil || out.RequestID == "" {
		t.Fatalf("unexpected delete response: %#v", out)
	}
}

func TestListContactGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/contacts/groups")
		if v := r.URL.Query().Get("limit"); v != "20" {
			t.Fatalf("limit = %q, want 20", v)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{
					"id":         "grp-1",
					"grant_id":   "abc-123",
					"group_type": "user",
					"name":       "My Group",
					"path":       "My Group",
				},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	lim := 20
	q := &models.ListContactGroupsQueryParams{Limit: &lim}
	out, err := c.Contacts().ListGroups(context.Background(), "abc-123", q)
	if err != nil {
		t.Fatalf("ListGroups error: %v", err)
	}
	if out == nil || out.RequestID == "" || len(out.Data) != 1 || out.Data[0].ID != "grp-1" {
		t.Fatalf("unexpected response: %#v", out)
	}
}
