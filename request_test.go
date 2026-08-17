package raiderio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/characters/profile" || r.URL.Query().Get("name") != "x" {
			t.Errorf("got %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"name":"X"}`))
	}))
	defer ts.Close()

	c := NewClient()
	c.ApiUrl = ts.URL
	got, err := getJSON[Character](c, context.Background(), "/characters/profile", url.Values{"name": {"x"}})
	if err != nil || got.Name != "X" {
		t.Fatalf("got %+v err %v", got, err)
	}
}
