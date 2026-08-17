package raiderio

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestGetMythicPlusAffixes(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "mythic_plus_affixes.json"))
	})

	got, err := c.GetMythicPlusAffixes(context.Background(), &AffixesQuery{Region: US, Locale: "en"})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/mythic-plus/affixes" {
		t.Fatalf("path %s", path)
	}
	if q.Get("region") != "us" || q.Get("locale") != "en" {
		t.Fatalf("query %v", q)
	}
	if got.Title == "" || len(got.AffixDetails) != 2 || got.AffixDetails[0].Name != "Xal'atath's Bargain: Pulsar" {
		t.Fatalf("got %+v", got)
	}
}

func TestGetMythicPlusAffixes_validate(t *testing.T) {
	c := mustNotHit(t)
	if _, err := c.GetMythicPlusAffixes(context.Background(), &AffixesQuery{}); !errors.Is(err, ErrInvalidRegion) {
		t.Fatalf("got %v want %v", err, ErrInvalidRegion)
	}
	if _, err := c.GetMythicPlusAffixes(context.Background(), nil); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("got %v want %v", err, ErrInvalidQuery)
	}
}

func TestGetMythicPlusStaticData(t *testing.T) {
	var path, exp string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		exp = r.URL.Query().Get("expansion_id")
		w.Write(testdata(t, "mythic_plus_static_data.json"))
	})

	got, err := c.GetMythicPlusStaticData(context.Background(), WAR_WITHIN)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/mythic-plus/static-data" || exp != "10" {
		t.Fatalf("path %s expansion_id %s", path, exp)
	}
	if len(got.Seasons) != 1 || got.Seasons[0].Slug != "season-tww-3" {
		t.Fatalf("got %+v", got)
	}
	if !got.Seasons[0].Starts.Us.Equal(time.Date(2025, 8, 12, 15, 0, 0, 0, time.UTC)) {
		t.Fatalf("starts %+v", got.Seasons[0].Starts)
	}
	if len(got.Seasons[0].Dungeons) != 1 || got.Seasons[0].Dungeons[0].Slug != "arakara-city-of-echoes" {
		t.Fatalf("dungeons %+v", got.Seasons[0].Dungeons)
	}
}

func TestGetMythicPlusRuns(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "mythic_plus_runs.json"))
	})

	got, err := c.GetMythicPlusRuns(context.Background(), &MythicPlusRunsQuery{
		Region: US, Season: "season-tww-3", Dungeon: "all", Affixes: "all", Page: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/mythic-plus/runs" {
		t.Fatalf("path %s", path)
	}
	if q.Get("region") != "us" || q.Get("season") != "season-tww-3" || q.Get("dungeon") != "all" ||
		q.Get("affixes") != "all" || q.Get("page") != "1" {
		t.Fatalf("query %v", q)
	}
	if len(got.Rankings) != 1 || got.Rankings[0].Rank != 1 || got.Rankings[0].Score != 551.4 {
		t.Fatalf("got %+v", got)
	}
	run := got.Rankings[0].Run
	if run.Dungeon.Slug != "ecodome-aldani" || run.MythicLevel != 24 {
		t.Fatalf("run %+v", run)
	}
	if len(run.WeeklyModifiers) != 1 || run.WeeklyModifiers[0].Name != "Tyrannical" {
		t.Fatalf("modifiers %+v", run.WeeklyModifiers)
	}
	if len(run.Roster) != 1 || run.Roster[0].Character.Name != "Silentdk" ||
		run.Roster[0].Character.Spec.Slug != "frost" || run.Roster[0].Character.Realm.Slug != "illidan" {
		t.Fatalf("roster %+v", run.Roster)
	}
}

func TestGetMythicPlusRuns_optional(t *testing.T) {
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		q = r.URL.Query()
		w.Write(testdata(t, "mythic_plus_runs.json"))
	})

	if _, err := c.GetMythicPlusRuns(context.Background(), &MythicPlusRunsQuery{}); err != nil {
		t.Fatal(err)
	}
	if len(q) != 0 {
		t.Fatalf("expected no query params, got %v", q)
	}
}

func TestGetMythicPlusRuns_validate(t *testing.T) {
	c := mustNotHit(t)
	if _, err := c.GetMythicPlusRuns(context.Background(), nil); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("nil query: got %v want %v", err, ErrInvalidQuery)
	}
}

func TestGetMythicPlusRunDetails(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "mythic_plus_run_details.json"))
	})

	got, err := c.GetMythicPlusRunDetails(context.Background(), &RunDetailsQuery{Season: "season-tww-3", ID: 41876796})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/mythic-plus/run-details" {
		t.Fatalf("path %s", path)
	}
	if q.Get("season") != "season-tww-3" || q.Get("id") != "41876796" {
		t.Fatalf("query %v", q)
	}
	if got.Score != 551.4 || got.KeystoneRunID != 41876796 || got.Dungeon.Name != "Eco-Dome Al'dani" {
		t.Fatalf("got %+v", got)
	}
	if !got.CompletedAt.Equal(time.Date(2026, 1, 21, 22, 29, 13, 0, time.UTC)) {
		t.Fatalf("completedAt %v", got.CompletedAt)
	}
}

func TestGetMythicPlusRunDetails_validate(t *testing.T) {
	c := mustNotHit(t)
	cases := []struct {
		q    RunDetailsQuery
		want error
	}{
		{RunDetailsQuery{ID: 1}, ErrInvalidSeason},
		{RunDetailsQuery{Season: "season-tww-3"}, ErrInvalidRunID},
		{RunDetailsQuery{Season: "season-tww-3", ID: -1}, ErrInvalidRunID},
	}
	for _, tc := range cases {
		_, err := c.GetMythicPlusRunDetails(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
	if _, err := c.GetMythicPlusRunDetails(context.Background(), nil); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("nil query: got %v want %v", err, ErrInvalidQuery)
	}
}

func TestGetMythicPlusScoreTiers(t *testing.T) {
	var season string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		season = r.URL.Query().Get("season")
		w.Write(testdata(t, "mythic_plus_score_tiers.json"))
	})

	got, err := c.GetMythicPlusScoreTiers(context.Background(), &ScoreTiersQuery{Season: "season-tww-3"})
	if err != nil {
		t.Fatal(err)
	}
	if season != "season-tww-3" {
		t.Fatalf("season %q", season)
	}
	if len(got) != 2 || got[0].Score != 4325 || got[0].RgbHex != "#ff8000" {
		t.Fatalf("got %+v", got)
	}
}

func TestGetMythicPlusScoreTiers_optional(t *testing.T) {
	var season string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		season = r.URL.Query().Get("season")
		w.Write(testdata(t, "mythic_plus_score_tiers.json"))
	})

	if _, err := c.GetMythicPlusScoreTiers(context.Background(), &ScoreTiersQuery{}); err != nil {
		t.Fatal(err)
	}
	if season != "" {
		t.Fatalf("season %q", season)
	}
}

func TestGetMythicPlusScoreTiers_validate(t *testing.T) {
	c := mustNotHit(t)
	if _, err := c.GetMythicPlusScoreTiers(context.Background(), nil); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("nil query: got %v want %v", err, ErrInvalidQuery)
	}
}

func TestGetMythicPlusSeasonCutoffs(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "season_cutoffs.json"))
	})

	got, err := c.GetMythicPlusSeasonCutoffs(context.Background(), &SeasonCutoffsQuery{Season: "season-tww-3", Region: US})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/mythic-plus/season-cutoffs" {
		t.Fatalf("path %s", path)
	}
	if q.Get("season") != "season-tww-3" || q.Get("region") != "us" {
		t.Fatalf("query %v", q)
	}
	if got.Cutoffs.Region.Slug != "us" || got.Cutoffs.P999 == nil {
		t.Fatalf("got %+v", got.Cutoffs)
	}
	if got.Cutoffs.P999.Horde.QuantileMinValue != 3731.37 || got.Cutoffs.P999.AllianceColor != "#f9763b" {
		t.Fatalf("p999 %+v", got.Cutoffs.P999)
	}
}

func TestGetMythicPlusSeasonCutoffs_validate(t *testing.T) {
	c := mustNotHit(t)
	cases := []struct {
		q    SeasonCutoffsQuery
		want error
	}{
		{SeasonCutoffsQuery{Region: US}, ErrInvalidSeason},
		{SeasonCutoffsQuery{Season: "season-tww-3"}, ErrInvalidRegion},
	}
	for _, tc := range cases {
		_, err := c.GetMythicPlusSeasonCutoffs(context.Background(), &tc.q)
		if !errors.Is(err, tc.want) {
			t.Fatalf("%+v: got %v want %v", tc.q, err, tc.want)
		}
	}
	if _, err := c.GetMythicPlusSeasonCutoffs(context.Background(), nil); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("nil query: got %v want %v", err, ErrInvalidQuery)
	}
}

func TestGetMythicPlusLeaderboardCapacity(t *testing.T) {
	var path string
	var q url.Values
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		q = r.URL.Query()
		w.Write(testdata(t, "leaderboard_capacity.json"))
	})

	got, err := c.GetMythicPlusLeaderboardCapacity(context.Background(), &LeaderboardCapacityQuery{Region: US, Realm: "illidan"})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/mythic-plus/leaderboard-capacity" {
		t.Fatalf("path %s", path)
	}
	if q.Get("region") != "us" || q.Get("realm") != "illidan" {
		t.Fatalf("query %v", q)
	}
	rl := got.RealmListing
	if rl.Region.Slug != "us" || len(rl.Affixes) != 1 || rl.Affixes[0].Name["en"] != "Xal'atath's Bargain: Pulsar" {
		t.Fatalf("got %+v", rl)
	}
	if len(rl.Realms) != 1 || rl.Realms[0].ConnectedRealms[0].Slug != "aegwynn" {
		t.Fatalf("realms %+v", rl.Realms)
	}
	if rl.Realms[0].Dungeons[0].Dungeon.Slug != "algethar-academy" {
		t.Fatalf("dungeons %+v", rl.Realms[0].Dungeons)
	}
}

func TestGetMythicPlusLeaderboardCapacity_validate(t *testing.T) {
	c := mustNotHit(t)
	if _, err := c.GetMythicPlusLeaderboardCapacity(context.Background(), &LeaderboardCapacityQuery{}); !errors.Is(err, ErrInvalidRegion) {
		t.Fatalf("got %v want %v", err, ErrInvalidRegion)
	}
	if _, err := c.GetMythicPlusLeaderboardCapacity(context.Background(), nil); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("got %v want %v", err, ErrInvalidQuery)
	}
}

func TestGetPeriods(t *testing.T) {
	var path string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Write(testdata(t, "periods.json"))
	})

	got, err := c.GetPeriods(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if path != "/periods" {
		t.Fatalf("path %s", path)
	}
	if len(got.Periods) != 2 || got.Periods[0].Region != "us" {
		t.Fatalf("got %+v", got)
	}
	if got.Periods[0].Current.Period != 1076 ||
		!got.Periods[0].Current.Start.Equal(time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)) {
		t.Fatalf("current %+v", got.Periods[0].Current)
	}
}
