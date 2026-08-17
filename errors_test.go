package raiderio

import "testing"

func TestWrapApiError(t *testing.T) {
	cases := []struct {
		msg  string
		want error
	}{
		{"Failed to find region usx", ErrInvalidRegion},
		{"Could not find requested character", ErrCharacterNotFound},
		{"Requested unsupported expansion_id", ErrUnsupportedExpac},
		{"Invalid request query input", ErrInvalidQuery},
		{"something else", ErrUnexpected},
	}
	for _, tc := range cases {
		got := wrapApiError(&apiErrorResponse{Message: tc.msg})
		if got != tc.want {
			t.Fatalf("%q: got %v want %v", tc.msg, got, tc.want)
		}
	}
}
