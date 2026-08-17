package raiderio_test

import (
	"encoding/json"
	"testing"

	"github.com/tmaffia/raiderio"
	"github.com/tmaffia/raiderio/regions"
)

func TestGuildRaidProgressionJSON(t *testing.T) {
	var g raiderio.Guild
	err := json.Unmarshal([]byte(`{"raid_progression":{"nerubar-palace":{"summary":"8/8 M"},"manaforge-omega":{"summary":"8/8 H"}}}`), &g)
	if err != nil {
		t.Fatal(err)
	}
	if g.RaidProgression["nerubar-palace"].Summary != "8/8 M" {
		t.Fatalf("got %#v", g.RaidProgression)
	}
}

func TestGetGuildRaidRankBySlug(t *testing.T) {
	testCases := []struct {
		region              *regions.Region
		realm               string
		name                string
		includeRandRankings bool
		raidSlug            string
		expectedErrMsg      string
	}{
		{region: regions.US, realm: "illidan", name: "warpath", raidSlug: "tier-mn-1", includeRandRankings: true},
		{region: regions.US, realm: "illidan", name: "warpath", raidSlug: "invalid raid slug", expectedErrMsg: "invalid raid", includeRandRankings: true},
		{region: regions.US, realm: "illidan", name: "warpath", raidSlug: "tier-mn-1",
			expectedErrMsg: "guild raid rankings field missing from api response", includeRandRankings: false},
	}

	for _, tc := range testCases {
		ctx, cancel := ctx()
		defer cancel()
		profile, err := c.GetGuild(ctx, &raiderio.GuildQuery{
			Region:       tc.region,
			Realm:        tc.realm,
			Name:         tc.name,
			RaidRankings: tc.includeRandRankings,
		})
		if err != nil {
			t.Fatalf("Error getting guild: %v", err)
		}

		rank, err := profile.GetGuildRaidRankBySlug(tc.raidSlug)
		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil && !(rank.Mythic.World > 0) {
			t.Fatalf("mythic guild ranking for raid: %v, got: %d, expected > 0",
				rank.RaidSlug, rank.Mythic.World)
		}
	}
}
