package raiderio

import (
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

func TestUnmarshalGuildBossKill(t *testing.T) {
	k, err := unmarshalGuildBossKill(testdata(t, "boss_kill.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !k.Kill.IsSuccess || k.Kill.Duration != 396990*time.Millisecond {
		t.Fatalf("kill %+v", k.Kill)
	}
	if len(k.Roster) != 1 {
		t.Fatalf("roster %d", len(k.Roster))
	}
	c := k.Roster[0]
	if c.Name != "Drbananaphd" || c.Class != "monk" || c.Spec != "windwalker" ||
		c.Realm != "illidan" || c.Region != "us" ||
		c.TalentLoadout.LoadoutText != "B0QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAESkkEtISEAAAAQikkAIJJJJpIpEiQASSSSaJBOAAA" ||
		c.Gear.ItemLevelEquipped != 401 {
		t.Fatalf("char %+v", c)
	}

	if _, err := unmarshalGuildBossKill([]byte(`{`)); err == nil {
		t.Fatal("expected error")
	}
}
