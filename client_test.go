package raiderio_test

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tmaffia/raiderio"
	"github.com/tmaffia/raiderio/expansions"
	"github.com/tmaffia/raiderio/regions"
)

var c *raiderio.Client
var defaultCtx context.Context
var cancel context.CancelFunc

func setup() {
	// Try to read .env and .env.local files
	for _, filename := range []string{".env", ".env.local", ".env.Local"} {
		func() {
			file, err := os.Open(filename)
			if err == nil {
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						key := strings.TrimSpace(parts[0])
						value := strings.TrimSpace(parts[1])
						os.Setenv(key, value)
					}
				}
			}
		}()
	}

	var opts []raiderio.ClientOption
	if key := os.Getenv("RAIDERIO_ACCESS_KEY"); key != "" {
		opts = append(opts, raiderio.WithAccessKey(key))
	}
	c = raiderio.NewClient(opts...)
	defaultCtx, cancel = context.WithTimeout(context.Background(), time.Second*30)
}

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNewClient(t *testing.T) {
	if c.ApiUrl != "https://raider.io/api/v1" {
		t.Errorf("NewClient apiUrl created incorrectly")
	}
}

func TestClient_WithAccessKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["access_key"]
		if !ok || len(keys[0]) < 1 {
			t.Errorf("access_key query param not found")
		}
		if keys[0] != "test_key" {
			t.Errorf("access_key expected: test_key, got: %v", keys[0])
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name": "Test Character"}`))
	}))
	defer ts.Close()

	client := raiderio.NewClient(raiderio.WithAccessKey("test_key"))
	client.ApiUrl = ts.URL

	_, err := client.GetCharacter(context.Background(), &raiderio.CharacterQuery{
		Region: regions.US,
		Realm:  "illidan",
		Name:   "test",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClient_NoAccessKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["access_key"]
		if ok && len(keys) > 0 {
			t.Errorf("access_key query param should NOT be present, got: %v", keys[0])
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name": "Test Character"}`))
	}))
	defer ts.Close()

	// Create a local client explicitly WITHOUT an AccessKey
	client := raiderio.NewClient()
	client.ApiUrl = ts.URL
	// Ensure AccessKey is empty (it should be by default, but being explicit)
	client.AccessKey = ""

	_, err := client.GetCharacter(context.Background(), &raiderio.CharacterQuery{
		Region: regions.US,
		Realm:  "illidan",
		Name:   "test",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetCharacterProfile(t *testing.T) {
	testCases := []struct {
		timeout        bool
		region         *regions.Region
		realm          string
		name           string
		expectedErrMsg string
		expectedName   string
	}{
		{region: regions.US, realm: "illidan", name: "highervalue", expectedName: "Highervalue"},
		{region: regions.US, realm: "", name: "highervalue", expectedErrMsg: "invalid realm"},
		{region: regions.US, realm: "illidan", name: "", expectedErrMsg: "invalid character name"},
		{region: nil, realm: "illidan", name: "highervalue", expectedErrMsg: "invalid region"},
		{region: &regions.Region{Slug: "badregion"}, realm: "illidan", name: "impossiblecharactername", expectedErrMsg: "invalid region"},
		{region: regions.US, realm: "illidan", name: "impossiblecharactername", expectedErrMsg: "character not found"},
		{region: regions.US, realm: "invalidrealm", name: "highervalue", expectedErrMsg: "invalid realm"},
		{timeout: true, region: regions.US, realm: "illidan", name: "highervalue", expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(defaultCtx, time.Millisecond*1)
			defer cancel()
		}

		profile, err := c.GetCharacter(ctx, &raiderio.CharacterQuery{
			Region: tc.region,
			Realm:  tc.realm,
			Name:   tc.name,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil && profile.Name != tc.expectedName {
			t.Fatalf("character name expected: %v, got: %v", tc.expectedName, profile.Name)
		}
	}
}

func TestGetCharacterWGear(t *testing.T) {
	testCases := []struct {
		timeout        bool
		region         *regions.Region
		realm          string
		name           string
		expectedErrMsg string
		expectedName   string
	}{
		{region: regions.US, realm: "illidan", name: "highervalue", expectedName: "Highervalue"},
		{timeout: true, region: regions.US, realm: "illidan", name: "highervalue", expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		profile, err := c.GetCharacter(ctx, &raiderio.CharacterQuery{
			Region: tc.region,
			Realm:  tc.realm,
			Name:   tc.name,
			Gear:   true,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil && profile.Name != tc.expectedName {
			t.Fatalf("character name expected: %v, got: %v. item level equipped: %d", tc.expectedName, profile.Name, profile.Gear.ItemLevelEquipped)
		}

		if err == nil && !(profile.Gear.ItemLevelEquipped > 0) {
			t.Fatalf("character item level equipped: %d, expected > 0", profile.Gear.ItemLevelEquipped)
		}
	}
}

func TestGetCharacterWTalents(t *testing.T) {
	cq := raiderio.CharacterQuery{
		Region:        regions.US,
		Realm:         "illidan",
		Name:          "highervalue",
		TalentLoadout: true,
	}

	profile, err := c.GetCharacter(defaultCtx, &cq)
	if err == nil && profile.TalentLoadout.LoadoutText == "" {
		t.Fatalf("talent loadout: %v expected to not be empty", profile.TalentLoadout.LoadoutText)
	}
}

func TestGetGuild(t *testing.T) {
	testCases := []struct {
		timeout        bool
		region         *regions.Region
		realm          string
		name           string
		expectedErrMsg string
		expectedName   string
	}{
		{region: regions.US, realm: "illidan", name: "warpath", expectedName: "Warpath"},
		{region: regions.US, realm: "", name: "warpath", expectedErrMsg: "invalid realm"},
		{region: regions.US, realm: "illidan", name: "", expectedErrMsg: "invalid guild name"},
		{region: nil, realm: "illidan", name: "highervalue", expectedErrMsg: "invalid region"},
		{region: &regions.Region{Slug: "badregion"}, realm: "illidan", name: "warpath", expectedErrMsg: "invalid region"},
		{region: regions.US, realm: "illidan", name: "impossible_guild_name", expectedErrMsg: "guild not found"},
		{region: regions.US, realm: "invalidrealm", name: "highervalue", expectedErrMsg: "invalid realm"},
		{timeout: true, region: regions.US, realm: "illidan", name: "highervalue", expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		profile, err := c.GetGuild(ctx, &raiderio.GuildQuery{
			Region: tc.region,
			Realm:  tc.realm,
			Name:   tc.name,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil && profile.Name != tc.expectedName {
			t.Fatalf("guild name expected: %v, got: %v.", tc.expectedName, profile.Name)
		}
	}
}

func TestGetGuildWMembers(t *testing.T) {
	testCases := []struct {
		region *regions.Region
		realm  string
		name   string
	}{
		{region: regions.US, realm: "illidan", name: "warpath"},
	}

	for range testCases {
		profile, err := c.GetGuild(defaultCtx, &raiderio.GuildQuery{
			Region:  regions.US,
			Realm:   "illidan",
			Name:    "warpath",
			Members: true,
		})

		if err != nil {
			t.Fatalf("Error getting guild: %v", err)
		}

		if !(len(profile.Members) > 0) {
			t.Fatalf("Error getting guild members")
		}
	}

}

func TestGetGuildWRaidProgression(t *testing.T) {
	testCases := []struct {
		region *regions.Region
		realm  string
		name   string
	}{
		{region: regions.US, realm: "illidan", name: "warpath"},
	}

	for range testCases {
		profile, err := c.GetGuild(defaultCtx, &raiderio.GuildQuery{
			Region:          regions.US,
			Realm:           "illidan",
			Name:            "warpath",
			RaidProgression: true,
		})

		if err != nil {
			t.Errorf("Error getting guild %v", err)
		}

		if profile.RaidProgression.NerubarPalace.Summary == "" {
			t.Fatalf("Error getting guild raid progression, got: %v", profile.RaidProgression.NerubarPalace.Summary)
		}
	}
}

func TestGetGuildWRaidRankings(t *testing.T) {
	testCases := []struct {
		timeout        bool
		region         *regions.Region
		realm          string
		name           string
		raidName       string
		expectedRank   int
		expectedErrMsg string
	}{
		{region: regions.US, realm: "illidan", name: "warpath",
			raidName: "nerubar-palace", expectedRank: 92},
		{timeout: true, region: regions.US, realm: "illidan", name: "warpath",
			raidName:       "nerubar-palace",
			expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		profile, err := c.GetGuild(ctx, &raiderio.GuildQuery{
			Region:       regions.US,
			Realm:        "illidan",
			Name:         "warpath",
			RaidRankings: true,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("error message got: %v, expected: %v",
				err.Error(), tc.expectedErrMsg)
		}

		if err == nil {
			rank := profile.RaidRankings[tc.raidName]
			if rank.Mythic.World != tc.expectedRank {
				t.Fatalf("mythic guild ranking for raid: %v, got: %d, expected: %d",
					rank.RaidSlug, rank.Mythic.World, tc.expectedRank)
			}
		}
	}
}

func TestGetGuildBossKill(t *testing.T) {
	testCases := []struct {
		region                *regions.Region
		realm                 string
		guildName             string
		raidSlug              string
		bossSlug              string
		difficulty            raiderio.RaidDifficulty
		expectedDefeatedAt    string
		expectedCharacterName string
		expectedErrMsg        string
		timeout               bool
	}{
		{region: regions.US, realm: "illidan", guildName: "warpath",
			raidSlug: "vault-of-the-incarnates", bossSlug: "terros",
			difficulty: raiderio.MYTHIC_RAID, expectedCharacterName: "Drbananaphd"},
		{region: nil, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			guildName: "warpath", raidSlug: "vault-of-the-incarnates",
			bossSlug: "terros", expectedErrMsg: "invalid region"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID,
			realm: "invalid-realm", guildName: "warpath", raidSlug: "vault-of-the-incarnates",
			bossSlug: "terros", expectedErrMsg: "invalid realm"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID,
			guildName: "warpath", raidSlug: "vault-of-the-incarnates",
			bossSlug: "terros", expectedErrMsg: "invalid realm"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			guildName: "impossible-guild_name", raidSlug: "vault-of-the-incarnates",
			bossSlug: "terros", expectedErrMsg: "guild not found"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			raidSlug: "vault-of-the-incarnates", bossSlug: "terros",
			expectedErrMsg: "invalid guild name"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			guildName: "warpath", raidSlug: "invalid-raid-slug", bossSlug: "terros",
			expectedErrMsg: "invalid raid"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			guildName: "warpath", bossSlug: "terros",
			expectedErrMsg: "invalid raid name"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			guildName: "warpath", raidSlug: "vault-of-the-incarnates",
			bossSlug: "invalid-boss-slug", expectedErrMsg: "invalid boss"},
		{region: regions.US, difficulty: raiderio.MYTHIC_RAID, realm: "illidan",
			guildName: "warpath", raidSlug: "vault-of-the-incarnates",
			expectedErrMsg: "invalid boss"},
		{region: regions.US, realm: "illidan", guildName: "warpath",
			raidSlug: "vault-of-the-incarnates", bossSlug: "terros",
			expectedErrMsg: "invalid raid difficulty"},
		{timeout: true, region: regions.US, realm: "illidan", guildName: "warpath",
			raidSlug: "vault-of-the-incarnates", bossSlug: "terros",
			difficulty:     raiderio.MYTHIC_RAID,
			expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		k, err := c.GetGuildBossKill(ctx, &raiderio.GuildBossKillQuery{
			Region:     tc.region,
			Realm:      tc.realm,
			GuildName:  tc.guildName,
			RaidSlug:   tc.raidSlug,
			BossSlug:   tc.bossSlug,
			Difficulty: tc.difficulty,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("error message got: %v, expected: %v", err.Error(), tc.expectedErrMsg)
		}

		if err == nil && !killIncludesCharacter(k, tc.expectedCharacterName) {
			t.Fatalf("boss kill character name expected: %v", tc.expectedCharacterName)
		}
	}
}

func TestGetRaids(t *testing.T) {
	testCases := []struct {
		timeout          bool
		expansion        expansions.Expansion
		raidName         string
		expectedRaidName string
		expectedErrMsg   string
	}{
		{expansion: expansions.DRAGONFLIGHT, raidName: "aberrus-the-shadowed-crucible", expectedRaidName: "Aberrus, the Shadowed Crucible"},
		{expansion: expansions.MIDNIGHT, raidName: "march-on-queldanas", expectedRaidName: "March on Quel'Danas"},
		{timeout: true, expansion: expansions.DRAGONFLIGHT, raidName: "aberrus-the-shadowed-crucible", expectedErrMsg: "raiderio api request timeout"},
		{expansion: 2, expectedErrMsg: "unsupported expansion"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		raids, err := c.GetRaids(ctx, tc.expansion)
		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil {
			raid, _ := raids.GetRaidBySlug(tc.raidName)
			if raid.Name != tc.expectedRaidName {
				t.Fatalf("expected raid name: %v, got: %v", tc.expectedRaidName, raid.Name)
			}
		}

	}
}

func TestGetRaidRankings(t *testing.T) {
	testCases := []struct {
		timeout        bool
		slug           string
		difficulty     raiderio.RaidDifficulty
		region         *regions.Region
		realm          string
		limit          int
		page           int
		expectedErrMsg string
	}{
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.WORLD, limit: 1},
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.US,
			realm: "proudmoore", limit: 1},
		{slug: "aberrus-the-shadowed-crucible", difficulty: "mythic", region: regions.EU, limit: 1},
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.US,
			realm: "illidan", limit: 1},
		{slug: "invalid raid slug", difficulty: raiderio.MYTHIC_RAID, region: regions.US, realm: "illidan",
			expectedErrMsg: "invalid raid"},
		{slug: "aberrus-the-shadowed-crucible", difficulty: "mythic", region: nil, realm: "illidan", expectedErrMsg: "invalid region"},
		{slug: "aberrus-the-shadowed-crucible", difficulty: "", region: regions.US, realm: "illidan",
			expectedErrMsg: "invalid raid difficulty"},
		{slug: "aberrus-the-shadowed-crucible", difficulty: "invalid-difficulty", region: regions.US, realm: "illidan",
			expectedErrMsg: "invalid raid difficulty"},
		{slug: "", difficulty: raiderio.MYTHIC_RAID, region: regions.US, realm: "illidan",
			expectedErrMsg: "invalid raid name"},
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.WORLD, limit: 20},
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.WORLD, limit: -20,
			expectedErrMsg: "limit must be a positive int"},
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.US, limit: 40, page: 2},
		{slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.US, limit: 40,
			page: -2, expectedErrMsg: "page must be a positive int"},
		{timeout: true, slug: "aberrus-the-shadowed-crucible", difficulty: raiderio.MYTHIC_RAID, region: regions.US,
			expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		rankings, err := c.GetRaidRankings(ctx, &raiderio.RaidQuery{
			Slug:       tc.slug,
			Difficulty: raiderio.RaidDifficulty(tc.difficulty),
			Region:     tc.region,
			Realm:      tc.realm,
			Limit:      tc.limit,
			Page:       tc.page,
		})
		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil && tc.limit != 0 {
			if len(rankings.RaidRanking) != tc.limit {
				t.Fatalf("expected results limit: %v, got: %v", tc.limit, len(rankings.RaidRanking))
			}
		} else if err == nil && len(rankings.RaidRanking) == 0 {
			t.Fatalf("expected results to not be empty")
		}
	}
}

// Tests if character is a part of the particular boss kill
func killIncludesCharacter(k *raiderio.BossKill, c string) bool {
	for _, v := range k.Roster {
		if v.Name == c {
			return true
		}
	}
	return false
}

func TestGetBossRankings(t *testing.T) {
	testCases := []struct {
		timeout        bool
		raidSlug       string
		bossSlug       string
		difficulty     raiderio.RaidDifficulty
		region         *regions.Region
		realm          string
		expectedErrMsg string
	}{
		{raidSlug: "nerubar-palace", bossSlug: "queen-ansurek", difficulty: raiderio.MYTHIC_RAID, region: regions.US},
		{raidSlug: "nerubar-palace", bossSlug: "queen-ansurek", difficulty: raiderio.MYTHIC_RAID, region: regions.US, realm: "illidan"},
		{raidSlug: "", bossSlug: "queen-ansurek", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "invalid raid name"},
		{raidSlug: "nerubar-palace", bossSlug: "", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "invalid boss"},
		{raidSlug: "nerubar-palace", bossSlug: "queen-ansurek", difficulty: "", region: regions.US, expectedErrMsg: "invalid raid difficulty"},
		{raidSlug: "nerubar-palace", bossSlug: "queen-ansurek", difficulty: raiderio.MYTHIC_RAID, region: nil, expectedErrMsg: "invalid region"},
		{timeout: true, raidSlug: "nerubar-palace", bossSlug: "queen-ansurek", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		rankings, err := c.GetBossRankings(ctx, &raiderio.BossRankingsQuery{
			RaidSlug:   tc.raidSlug,
			BossSlug:   tc.bossSlug,
			Difficulty: tc.difficulty,
			Region:     tc.region,
			Realm:      tc.realm,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil {
			if len(rankings.BossRankings) == 0 {
				t.Fatalf("expected boss rankings to not be empty")
			}
		}
	}
}

func TestGetHallOfFame(t *testing.T) {
	// Create a client without access key
	noAuthClient := raiderio.NewClient()
	noAuthClient.AccessKey = ""

	testCases := []struct {
		timeout        bool
		raidSlug       string
		difficulty     raiderio.RaidDifficulty
		region         *regions.Region
		expectedErrMsg string
	}{
		{raidSlug: "nerubar-palace", difficulty: raiderio.MYTHIC_RAID, region: regions.US},
		{raidSlug: "", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "invalid raid name"},
		{raidSlug: "nerubar-palace", difficulty: "", region: regions.US, expectedErrMsg: "invalid raid difficulty"},
		{raidSlug: "nerubar-palace", difficulty: raiderio.MYTHIC_RAID, region: nil, expectedErrMsg: "invalid region"},
		{timeout: true, raidSlug: "nerubar-palace", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		hof, err := noAuthClient.GetHallOfFame(ctx, &raiderio.HallOfFameQuery{
			RaidSlug:   tc.raidSlug,
			Difficulty: tc.difficulty,
			Region:     tc.region,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil {
			if len(hof.HallOfFame.BossKills) == 0 {
				t.Fatalf("expected hall of fame to not be empty")
			}
		}
	}
}

func TestGetRaidProgression(t *testing.T) {
	// Create a client without access key to test if that's the issue
	noAuthClient := raiderio.NewClient()
	noAuthClient.AccessKey = ""

	testCases := []struct {
		timeout        bool
		raidSlug       string
		difficulty     raiderio.RaidDifficulty
		region         *regions.Region
		expectedErrMsg string
	}{
		{raidSlug: "nerubar-palace", difficulty: raiderio.MYTHIC_RAID, region: regions.US},
		{raidSlug: "", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "invalid raid name"},
		{raidSlug: "nerubar-palace", difficulty: "", region: regions.US, expectedErrMsg: "invalid raid difficulty"},
		{raidSlug: "nerubar-palace", difficulty: raiderio.MYTHIC_RAID, region: nil, expectedErrMsg: "invalid region"},
		{timeout: true, raidSlug: "nerubar-palace", difficulty: raiderio.MYTHIC_RAID, region: regions.US, expectedErrMsg: "raiderio api request timeout"},
	}

	for _, tc := range testCases {
		ctx := defaultCtx
		var cancel context.CancelFunc
		if tc.timeout {
			ctx, cancel = context.WithTimeout(ctx, time.Millisecond*1)
			defer cancel()
		}

		prog, err := noAuthClient.GetRaidProgression(ctx, &raiderio.RaidProgressionQuery{
			RaidSlug:   tc.raidSlug,
			Difficulty: tc.difficulty,
			Region:     tc.region,
		})

		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Fatalf("expected error: %v, got: %v", tc.expectedErrMsg, err.Error())
		}

		if err == nil {
			if len(prog.Progression) == 0 {
				t.Fatalf("expected raid progression to not be empty")
			}
		}
	}
}
