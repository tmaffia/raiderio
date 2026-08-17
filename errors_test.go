package raiderio

import (
	"context"
	"errors"
	"testing"
)

type hiddenErr struct{ error }

func (e hiddenErr) Error() string { return "transport failed" }
func (e hiddenErr) Unwrap() error { return e.error }

func TestWrapApiError(t *testing.T) {
	cases := []struct {
		msg  string
		want error
	}{
		{"Failed to find region usx", ErrInvalidRegion},
		{"Failed to find realm foo", ErrInvalidRealm},
		{"Failed to find raid bar", ErrInvalidRaid},
		{"Failed to find boss baz", ErrInvalidBoss},
		{"Could not find requested character", ErrCharacterNotFound},
		{"Could not find requested guild", ErrGuildNotFound},
		{"Could not find requested raid", ErrInvalidRaid},
		{"Requested unsupported expansion_id", ErrUnsupportedExpac},
		{"Invalid request query input", ErrInvalidQuery},
		{`"region" must be one of [us, eu, tw, kr, cn, world]`, ErrInvalidRegion},
		{"Could not find data for season bogus", ErrInvalidSeason},
		{"could not find keystone run", ErrRunNotFound},
		{"something else", ErrUnexpected},
	}
	for _, tc := range cases {
		got := wrapApiError(&apiErrorResponse{Message: tc.msg})
		if !errors.Is(got, tc.want) {
			t.Fatalf("%q: got %v want %v", tc.msg, got, tc.want)
		}
	}
}

func TestWrapHttpError(t *testing.T) {
	if got := wrapHttpError(hiddenErr{context.DeadlineExceeded}); !errors.Is(got, ErrApiTimeout) {
		t.Fatal(got)
	}
	if got := wrapHttpError(hiddenErr{context.Canceled}); !errors.Is(got, ErrApiTimeout) {
		t.Fatal(got)
	}
	if got := wrapHttpError(errors.New("connection refused")); !errors.Is(got, ErrUnexpected) {
		t.Fatal(got)
	}
}
