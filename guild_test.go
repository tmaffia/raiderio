package raiderio_test

import (
	"encoding/json"
	"testing"

	"github.com/tmaffia/raiderio"
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
