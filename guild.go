package raiderio

import "time"

// GuildQuery is a struct that represents the query parameters
// sent for a guild profile request
// Supports optional request fields: members, raid_progression, raid_rankings
type GuildQuery struct {
	Region          Region
	Realm           string
	Name            string
	Members         bool
	RaidProgression bool
	RaidRankings    bool
}

// Guild is a struct that represents the response from
// a guild profile request
type Guild struct {
	Name            string                      `json:"name"`
	Faction         string                      `json:"faction"`
	Region          string                      `json:"region"`
	Realm           string                      `json:"realm"`
	LastCrawledAt   time.Time                   `json:"last_crawled_at"`
	ProfileUrl      string                      `json:"profile_url"`
	Members         []Member                    `json:"members"`
	RaidProgression map[string]RaidProgression  `json:"raid_progression"`
	RaidRankings    map[string]GuildRaidRanking `json:"raid_rankings"`
}

// Member is a struct that represents a member of a guild
// in a guild profile response
type Member struct {
	Rank      int       `json:"rank"`
	Character Character `json:"character"`
}

// createGuildQuery creates and validates a GuildQuery struct
// It returns an error if any of the required parameters are empty
// or if the fields are invalid
func createGuildQuery(gq *GuildQuery) ([]string, error) {
	if gq == nil {
		return nil, ErrInvalidQuery
	}

	if gq.Region == "" {
		return nil, ErrInvalidRegion
	}

	if gq.Realm == "" {
		return nil, ErrInvalidRealm
	}

	if gq.Name == "" {
		return nil, ErrInvalidGuildName
	}

	var fields []string
	if gq.Members {
		fields = append(fields, "members")
	}

	if gq.RaidProgression {
		fields = append(fields, "raid_progression")
	}

	if gq.RaidRankings {
		fields = append(fields, "raid_rankings")
	}
	return fields, nil
}
