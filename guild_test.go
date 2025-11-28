package raiderio_test

import (
	"testing"

	"github.com/tmaffia/raiderio"
	"github.com/tmaffia/raiderio/regions"
)

func TestGetGuildRaidRankBySlug(t *testing.T) {
	testCases := []struct {
		region              *regions.Region
		realm               string
		name                string
		includeRandRankings bool
		raidSlug            string
		expectedErrMsg      string
	}{
		{region: regions.US, realm: "illidan", name: "warpath", raidSlug: "nerubar-palace", includeRandRankings: true},
		{region: regions.US, realm: "illidan", name: "warpath", raidSlug: "invalid raid slug", expectedErrMsg: "invalid raid", includeRandRankings: true},
		{region: regions.US, realm: "illidan", name: "warpath", raidSlug: "nerubar-palace",
			expectedErrMsg: "guild raid rankings field missing from api response", includeRandRankings: false},
	}

	for _, tc := range testCases {
		profile, err := c.GetGuild(defaultCtx, &raiderio.GuildQuery{
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
