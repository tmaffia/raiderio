package raiderio

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGetJSON(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/characters/profile" || r.URL.Query().Get("name") != "x" {
			t.Errorf("got %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		w.Write([]byte(`{"name":"X"}`))
	})
	got, err := getJSON[Character](c, context.Background(), "/characters/profile", url.Values{"name": {"x"}})
	if err != nil || got.Name != "X" {
		t.Fatalf("got %+v err %v", got, err)
	}
}

func TestGetJSON_apiError(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Could not find requested character"}`))
	})
	_, err := getJSON[Character](c, context.Background(), "/x", url.Values{})
	if !errors.Is(err, ErrCharacterNotFound) {
		t.Fatal(err)
	}
}

func TestGetJSON_badJSON(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{`))
	})
	_, err := getJSON[Character](c, context.Background(), "/x", url.Values{})
	if err == nil || err.Error() != "error unmarshalling response" {
		t.Fatal(err)
	}
}

func TestGetJSON_timeout(t *testing.T) {
	c := testClient(t, func(http.ResponseWriter, *http.Request) {})
	ctx, cancel := context.WithDeadline(context.Background(), time.Time{})
	defer cancel()
	_, err := getJSON[Character](c, ctx, "/x", url.Values{})
	if !errors.Is(err, ErrApiTimeout) {
		t.Fatal(err)
	}
}

func TestGetJSON_accessKey(t *testing.T) {
	var key string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key = r.URL.Query().Get("access_key")
		w.Write([]byte(`{"name":"X"}`))
	}))
	t.Cleanup(ts.Close)

	c := NewClient("test_key")
	c.ApiUrl = ts.URL
	if _, err := getJSON[Character](c, context.Background(), "/x", url.Values{}); err != nil {
		t.Fatal(err)
	}
	if key != "test_key" {
		t.Fatalf("access_key %q", key)
	}

	c = NewClient()
	c.ApiUrl = ts.URL
	key = "sentinel"
	if _, err := getJSON[Character](c, context.Background(), "/x", url.Values{}); err != nil {
		t.Fatal(err)
	}
	if key != "" {
		t.Fatalf("access_key %q", key)
	}
}
