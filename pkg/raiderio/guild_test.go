package raiderio

import "testing"

func TestCreateGuildQuery(t *testing.T) {
	gq := GuildQuery{
		Region: "us",
		Realm:  "illidan",
		Name:   "liquid",
	}

	err := createGuildQuery(&gq)
	if err != nil {
		t.Errorf("Error creating guild query")
	}
	t.Logf("%+v", gq)
}

func TestCreateGuildQueryWMembers(t *testing.T) {
	gq := GuildQuery{
		Region:  "us",
		Realm:   "illidan",
		Name:    "liquid",
		Members: true,
	}

	err := createGuildQuery(&gq)
	if err != nil {
		t.Errorf("Error creating guild query")
	}
	if gq.fields[0] != "members" {
		t.Errorf("Error creating guild query")
	}
	t.Logf("%+v", gq)
}

func TestCreateGuildQueryWRaidProgression(t *testing.T) {
	gq := GuildQuery{
		Region:          "us",
		Realm:           "illidan",
		Name:            "liquid",
		RaidProgression: true,
	}

	err := createGuildQuery(&gq)
	if err != nil {
		t.Errorf("Error creating guild query")
	}
	if gq.fields[0] != "raid_progression" {
		t.Errorf("Error creating guild query")
	}
	t.Logf("%+v", gq)
}
