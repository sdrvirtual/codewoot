package codechat

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteInstanceByName_UsesFetchedInstanceToken(t *testing.T) {
	var deleteAuth string
	var deletePath string
	var deleteQuery string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/instance/fetchInstances":
			if r.Header.Get("apikey") != "global-token" {
				t.Fatalf("expected global apikey, got %q", r.Header.Get("apikey"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[
				{"id":1,"name":"other","Auth":{"token":"other-token"}},
				{"id":2,"name":"target","Auth":{"token":"target-token"}}
			]`))
		case "/instance/delete/target":
			deleteAuth = r.Header.Get("Authorization")
			deletePath = r.URL.Path
			deleteQuery = r.URL.RawQuery
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	c, err := New(ts.URL, "global-token")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	_, err = c.DeleteInstanceByName(t.Context(), "target")
	if err != nil {
		t.Fatalf("DeleteInstanceByName() error: %v", err)
	}
	if deletePath != "/instance/delete/target" {
		t.Fatalf("unexpected delete path: %s", deletePath)
	}
	if deleteQuery != "force=true" {
		t.Fatalf("expected force=true query, got %q", deleteQuery)
	}
	if deleteAuth != "Bearer target-token" {
		t.Fatalf("expected fetched instance token, got %q", deleteAuth)
	}
}

func TestDeleteInstanceByName_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instance/fetchInstances" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	c, err := New(ts.URL, "global-token")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	_, err = c.DeleteInstanceByName(t.Context(), "missing")
	if err == nil {
		t.Fatalf("expected not found error")
	}
}
