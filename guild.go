package raiderio

import (
	"net/url"
	"strings"
	"time"
)

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

// validate checks the required fields of a GuildQuery.
func (q *GuildQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}

	if q.Region == "" {
		return ErrInvalidRegion
	}

	if q.Realm == "" {
		return ErrInvalidRealm
	}

	if q.Name == "" {
		return ErrInvalidGuildName
	}

	return nil
}

// params builds the URL parameters for a guild profile request.
func (q *GuildQuery) params() url.Values {
	params := url.Values{
		"region": {string(q.Region)},
		"realm":  {q.Realm},
		"name":   {q.Name},
	}
	if fields := q.fields(); len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	return params
}

// fields builds the optional request fields from the boolean toggles.
func (q *GuildQuery) fields() []string {
	var fields []string
	if q.Members {
		fields = append(fields, "members")
	}

	if q.RaidProgression {
		fields = append(fields, "raid_progression")
	}

	if q.RaidRankings {
		fields = append(fields, "raid_rankings")
	}
	return fields
}
