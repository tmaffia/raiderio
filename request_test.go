package raiderio

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestGetJSONFor(t *testing.T) {
	var path, region string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		region = r.URL.Query().Get("region")
		w.Write([]byte(`{"name":"X"}`))
	})

	got, err := getJSONFor[Character](c, context.Background(), "/characters/profile", &CharacterQuery{
		Region: US, Realm: "illidan", Name: "x",
	})
	if err != nil || got.Name != "X" {
		t.Fatalf("got %+v err %v", got, err)
	}
	if path != "/characters/profile" || region != "us" {
		t.Fatalf("path %s region %q", path, region)
	}
}

func TestGetJSONFor_nilQuery(t *testing.T) {
	c := mustNotHit(t)
	if _, err := getJSONFor[Character](c, context.Background(), "/x", (*CharacterQuery)(nil)); !errors.Is(err, ErrInvalidQuery) {
		t.Fatal(err)
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
	if err.Error() != "character not found: Could not find requested character" {
		t.Fatal(err)
	}
}

func TestGetJSON_unexpectedPreservesStatusAndBody(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"backend exploded"}`))
	})
	_, err := getJSON[Character](c, context.Background(), "/x", url.Values{})
	if !errors.Is(err, ErrUnexpected) {
		t.Fatal(err)
	}
	if err.Error() != `unexpected error: 500 backend exploded` {
		t.Fatal(err)
	}
}

func TestGetJSON_badJSON(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{`))
	})
	_, err := getJSON[Character](c, context.Background(), "/x", url.Values{})
	if err == nil || !strings.Contains(err.Error(), "error unmarshalling response") {
		t.Fatal(err)
	}
}

func TestGetJSON_timeout(t *testing.T) {
	c := testClient(t, func(http.ResponseWriter, *http.Request) {})
	ctx, cancel := context.WithDeadline(context.Background(), time.Time{})
	defer cancel()
	_, err := getJSON[Character](c, ctx, "/x", url.Values{})
	if !errors.Is(err, ErrApiTimeout) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
}

func TestGetJSON_canceled(t *testing.T) {
	c := testClient(t, func(http.ResponseWriter, *http.Request) {})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := getJSON[Character](c, ctx, "/x", url.Values{})
	if !errors.Is(err, ErrRequestCanceled) || !errors.Is(err, context.Canceled) {
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
