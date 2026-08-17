package raiderio

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)
	c := NewClient()
	c.ApiUrl = ts.URL
	return c
}

func mustNotHit(t *testing.T) *Client {
	t.Helper()
	return testClient(t, func(http.ResponseWriter, *http.Request) {
		t.Fatal("request should not be sent")
	})
}

func TestNewClient(t *testing.T) {
	if NewClient().ApiUrl != "https://raider.io/api/v1" {
		t.Fatal(NewClient().ApiUrl)
	}
	if NewClient().AccessKey != "" {
		t.Fatal("expected empty key")
	}
	if NewClient("k").AccessKey != "k" {
		t.Fatal("expected key")
	}
}

func TestGetCharacter(t *testing.T) {
	var path, fields string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		fields = r.URL.Query().Get("fields")
		q := r.URL.Query()
		if q.Get("region") != "us" || q.Get("realm") != "illidan" || q.Get("name") != "highervalue" {
			t.Errorf("query %v", q)
		}
		if fields == "" {
			w.Write(testdata(t, "character.json"))
			return
		}
		w.Write(testdata(t, "character_fields.json"))
	})

	p, err := c.GetCharacter(context.Background(), &CharacterQuery{
		Region: US, Realm: "illidan", Name: "highervalue",
	})
	if err != nil || p.Name != "Highervalue" || p.ActiveSpec != "Fury" || p.ActiveRole != "DPS" {
		t.Fatalf("got %+v err %v", p, err)
	}
	if path != "/characters/profile" {
		t.Fatalf("path %s", path)
	}

	p, err = c.GetCharacter(context.Background(), &CharacterQuery{
		Region: US, Realm: "illidan", Name: "highervalue",
		Gear: true, TalentLoadout: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if fields != "talents,gear" {
		t.Fatalf("fields %q", fields)
	}
	if p.Gear.ItemLevelEquipped <= 0 || p.TalentLoadout.LoadoutText == "" {
		t.Fatalf("got %+v", p)
	}
}

func TestGetCharacter_validate(t *testing.T) {
	c := mustNotHit(t)
	cases := []struct {
		q    CharacterQuery
		want error
	}{
		{CharacterQuery{Realm: "illidan", Name: "x"}, ErrInvalidRegion},
		{CharacterQuery{Region: US, Name: "x"}, ErrInvalidRealm},
		{CharacterQuery{Region: US, Realm: "illidan"}, ErrInvalidCharName},
	}
	for _, tc := range cases {
		_, err := c.GetCharacter(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
}

func TestGetGuild(t *testing.T) {
	var path, fields string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		fields = r.URL.Query().Get("fields")
		w.Write(testdata(t, "guild.json"))
	})

	g, err := c.GetGuild(context.Background(), &GuildQuery{
		Region: US, Realm: "illidan", Name: "warpath",
		Members: true, RaidProgression: true, RaidRankings: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/guilds/profile" {
		t.Fatalf("path %s", path)
	}
	if fields != "members,raid_progression,raid_rankings" {
		t.Fatalf("fields %q", fields)
	}
	if g.Name != "Warpath" || len(g.Members) == 0 {
		t.Fatalf("got %+v", g)
	}
	if g.RaidProgression["tier-mn-1"].Summary != "9/9 M" {
		t.Fatalf("progression %+v", g.RaidProgression)
	}
	if g.RaidRankings["tier-mn-1"].Mythic.World == 0 {
		t.Fatalf("rankings %+v", g.RaidRankings)
	}
}

func TestGetGuild_validate(t *testing.T) {
	c := mustNotHit(t)
	cases := []struct {
		q    GuildQuery
		want error
	}{
		{GuildQuery{Realm: "illidan", Name: "x"}, ErrInvalidRegion},
		{GuildQuery{Region: US, Name: "x"}, ErrInvalidRealm},
		{GuildQuery{Region: US, Realm: "illidan"}, ErrInvalidGuildName},
	}
	for _, tc := range cases {
		_, err := c.GetGuild(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
}

func TestGetRaids(t *testing.T) {
	var path, exp string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		exp = r.URL.Query().Get("expansion_id")
		w.Write(testdata(t, "raids.json"))
	})

	raids, err := c.GetRaids(context.Background(), WAR_WITHIN)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/raiding/static-data" || exp != "10" {
		t.Fatalf("path %s expansion_id %s", path, exp)
	}
	r, err := raids.GetRaidBySlug("nerubar-palace")
	if err != nil || r.Name != "Nerub-ar Palace" {
		t.Fatalf("got %+v err %v", r, err)
	}
}

func TestGetRaidRankings(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "raid_rankings.json"))
	})

	got, err := c.GetRaidRankings(context.Background(), &RaidQuery{
		Slug: "nerubar-palace", Difficulty: MYTHIC_RAID, Region: US,
		Realm: "illidan", Limit: 1, Page: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/raiding/raid-rankings" {
		t.Fatalf("path %s", path)
	}
	if q.Get("raid") != "nerubar-palace" || q.Get("difficulty") != "mythic" ||
		q.Get("region") != "us" || q.Get("realm") != "illidan" ||
		q.Get("limit") != "1" || q.Get("page") != "2" {
		t.Fatalf("query %v", q)
	}
	if len(got.RaidRanking) != 1 || got.RaidRanking[0].Guild.Name != "Liquid" {
		t.Fatalf("got %+v", got)
	}
}

func TestGetRaidRankings_validate(t *testing.T) {
	c := mustNotHit(t)
	ok := RaidQuery{Slug: "x", Difficulty: MYTHIC_RAID, Region: US}
	cases := []struct {
		q    RaidQuery
		want error
	}{
		{RaidQuery{Difficulty: MYTHIC_RAID, Region: US}, ErrInvalidRaidName},
		{RaidQuery{Slug: "x", Region: US}, ErrInvalidRaidDiff},
		{RaidQuery{Slug: "x", Difficulty: "nope", Region: US}, ErrInvalidRaidDiff},
		{RaidQuery{Slug: "x", Difficulty: MYTHIC_RAID}, ErrInvalidRegion},
		{func() RaidQuery { q := ok; q.Limit = -1; return q }(), ErrLimitOutOfBounds},
		{func() RaidQuery { q := ok; q.Page = -1; return q }(), ErrPageOutOfBounds},
	}
	for _, tc := range cases {
		_, err := c.GetRaidRankings(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
}

func TestGetGuildBossKill(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "boss_kill.json"))
	})

	k, err := c.GetGuildBossKill(context.Background(), &GuildBossKillQuery{
		Region: US, Realm: "illidan", GuildName: "warpath",
		RaidSlug: "vault-of-the-incarnates", BossSlug: "terros", Difficulty: MYTHIC_RAID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/guilds/boss-kill" {
		t.Fatalf("path %s", path)
	}
	if q.Get("raid") != "vault-of-the-incarnates" || q.Get("boss") != "terros" ||
		q.Get("guild") != "warpath" || q.Get("difficulty") != "mythic" {
		t.Fatalf("query %v", q)
	}
	if len(k.Roster) == 0 || k.Roster[0].Name != "Drbananaphd" {
		t.Fatalf("got %+v", k)
	}
	if k.Kill.Duration != 396990*time.Millisecond || !k.Kill.IsSuccess {
		t.Fatalf("kill %+v", k.Kill)
	}
}

func TestGetGuildBossKill_validate(t *testing.T) {
	c := mustNotHit(t)
	ok := GuildBossKillQuery{
		Region: US, Realm: "illidan", GuildName: "warpath",
		RaidSlug: "x", BossSlug: "y", Difficulty: MYTHIC_RAID,
	}
	cases := []struct {
		mut  func(*GuildBossKillQuery)
		want error
	}{
		{func(q *GuildBossKillQuery) { q.Region = "" }, ErrInvalidRegion},
		{func(q *GuildBossKillQuery) { q.Realm = "" }, ErrInvalidRealm},
		{func(q *GuildBossKillQuery) { q.GuildName = "" }, ErrInvalidGuildName},
		{func(q *GuildBossKillQuery) { q.RaidSlug = "" }, ErrInvalidRaidName},
		{func(q *GuildBossKillQuery) { q.BossSlug = "" }, ErrInvalidBoss},
		{func(q *GuildBossKillQuery) { q.Difficulty = "" }, ErrInvalidRaidDiff},
	}
	for _, tc := range cases {
		q := ok
		tc.mut(&q)
		_, err := c.GetGuildBossKill(context.Background(), &q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", q, err, tc.want)
		}
	}
}

func TestGetBossRankings(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "boss_rankings.json"))
	})

	got, err := c.GetBossRankings(context.Background(), &BossRankingsQuery{
		RaidSlug: "nerubar-palace", BossSlug: "queen-ansurek",
		Difficulty: MYTHIC_RAID, Region: US, Realm: "illidan",
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/raiding/boss-rankings" {
		t.Fatalf("path %s", path)
	}
	if q.Get("raid") != "nerubar-palace" || q.Get("boss") != "queen-ansurek" || q.Get("realm") != "illidan" {
		t.Fatalf("query %v", q)
	}
	if len(got.BossRankings) == 0 {
		t.Fatal("empty rankings")
	}
}

func TestGetBossRankings_validate(t *testing.T) {
	c := mustNotHit(t)
	cases := []struct {
		q    BossRankingsQuery
		want error
	}{
		{BossRankingsQuery{BossSlug: "x", Difficulty: MYTHIC_RAID, Region: US}, ErrInvalidRaidName},
		{BossRankingsQuery{RaidSlug: "x", Difficulty: MYTHIC_RAID, Region: US}, ErrInvalidBoss},
		{BossRankingsQuery{RaidSlug: "x", BossSlug: "y", Region: US}, ErrInvalidRaidDiff},
		{BossRankingsQuery{RaidSlug: "x", BossSlug: "y", Difficulty: MYTHIC_RAID}, ErrInvalidRegion},
	}
	for _, tc := range cases {
		_, err := c.GetBossRankings(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
}

func TestGetHallOfFame(t *testing.T) {
	var path string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Write(testdata(t, "hall_of_fame.json"))
	})

	got, err := c.GetHallOfFame(context.Background(), &RaidQuery{
		Slug: "nerubar-palace", Difficulty: MYTHIC_RAID, Region: US,
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/raiding/hall-of-fame" {
		t.Fatalf("path %s", path)
	}
	if len(got.HallOfFame.BossKills) == 0 {
		t.Fatal("empty hall of fame")
	}
}

func TestGetRaidProgression(t *testing.T) {
	var path string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Write(testdata(t, "raid_progression.json"))
	})

	got, err := c.GetRaidProgression(context.Background(), &RaidQuery{
		Slug: "nerubar-palace", Difficulty: MYTHIC_RAID, Region: US,
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/raiding/progression" {
		t.Fatalf("path %s", path)
	}
	if len(got.Progression) == 0 {
		t.Fatal("empty progression")
	}
}

func TestValidateRaidQuery(t *testing.T) {
	c := mustNotHit(t)
	cases := []struct {
		q    RaidQuery
		want error
	}{
		{RaidQuery{Difficulty: MYTHIC_RAID, Region: US}, ErrInvalidRaidName},
		{RaidQuery{Slug: "x", Region: US}, ErrInvalidRaidDiff},
		{RaidQuery{Slug: "x", Difficulty: MYTHIC_RAID}, ErrInvalidRegion},
	}
	for _, tc := range cases {
		_, err := c.GetHallOfFame(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
}
