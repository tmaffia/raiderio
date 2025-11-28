package raiderio_test

import (
	"testing"

	"github.com/tmaffia/raiderio/expansions"
)

func TestGetRaidBySlug(t *testing.T) {
	testCases := []struct {
		slug           string
		expectedName   string
		expectedErrMsg string
	}{
		{slug: "nerubar-palace", expectedName: "Nerub-ar Palace"},
		{slug: "invalid raid slug", expectedErrMsg: "invalid raid"},
		{slug: "nerubar-palaceinvalid raid slug", expectedErrMsg: "invalid raid"},
	}

	raids, err := c.GetRaids(defaultCtx, expansions.WAR_WITHIN)
	if err != nil {
		t.Fatalf("Error getting raids: %v", err)
	}

	for _, tc := range testCases {
		raid, err := raids.GetRaidBySlug(tc.slug)
		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil && raid.Name != tc.expectedName {
			t.Fatalf("expected raid name: %v, got: %v", tc.expectedName, raid.Name)
		}
	}
}
