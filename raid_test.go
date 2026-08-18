package raiderio

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestGetRaidBySlug(t *testing.T) {
	raids := &Raids{Raids: []Raid{{Slug: "nerubar-palace", Name: "Nerub-ar Palace"}}}
	r, err := raids.GetRaidBySlug("nerubar-palace")
	if err != nil || r.Name != "Nerub-ar Palace" {
		t.Fatalf("got %+v err %v", r, err)
	}
	if _, err := raids.GetRaidBySlug("nope"); !errors.Is(err, ErrInvalidRaid) {
		t.Fatal(err)
	}
}

func TestMapBossKill(t *testing.T) {
	var resp bossKillResp
	if err := json.Unmarshal(testdata(t, "boss_kill.json"), &resp); err != nil {
		t.Fatal(err)
	}
	k := mapBossKill(resp)
	if !k.Kill.IsSuccess || k.Kill.Duration != 396990*time.Millisecond {
		t.Fatalf("kill %+v", k.Kill)
	}
	if !k.Kill.PulledAt.Equal(time.Date(2022, 12, 27, 5, 17, 25, 488000000, time.UTC)) ||
		!k.Kill.DefeatedAt.Equal(time.Date(2022, 12, 27, 5, 24, 2, 478000000, time.UTC)) {
		t.Fatalf("kill times %+v", k.Kill)
	}
	if len(k.Roster) != 1 {
		t.Fatalf("roster %d", len(k.Roster))
	}
	c := k.Roster[0]
	if c.Name != "Drbananaphd" || c.Class != "monk" || c.Spec != "windwalker" ||
		c.Realm != "illidan" || c.Region != "us" ||
		c.TalentLoadout.LoadoutSpecID != 269 ||
		c.TalentLoadout.LoadoutText != "B0QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAESkkEtISEAAAAQikkAIJJJJpIpEiQASSSSaJBOAAA" ||
		c.Gear.ItemLevelEquipped != 401 {
		t.Fatalf("char %+v", c)
	}
}
