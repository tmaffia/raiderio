package raiderio

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Base URL for the Raider.IO API
const baseUrl string = "https://raider.io/api"

// Client is the main struct for interacting with the Raider.IO API
type Client struct {
	ApiUrl     string
	AccessKey  string
	HttpClient *http.Client
}

// NewClient creates a Client. Pass an access key for authenticated requests.
func NewClient(accessKey ...string) *Client {
	c := &Client{
		ApiUrl:     baseUrl + "/v1",
		HttpClient: &http.Client{},
	}
	if len(accessKey) > 0 {
		c.AccessKey = accessKey[0]
	}
	return c
}

// GetCharacter retrieves a character profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Character struct
func (c *Client) GetCharacter(ctx context.Context, q *CharacterQuery) (*Character, error) {
	return getJSONFor[Character](c, ctx, "/characters/profile", q)
}

// GetGuild retrieves a guild profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Guild struct
func (c *Client) GetGuild(ctx context.Context, q *GuildQuery) (*Guild, error) {
	return getJSONFor[Guild](c, ctx, "/guilds/profile", q)
}

// GetRaids retrieves a list of raids from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Raids struct
// Takes an Expansion enum as a parameter, in addition to context.Context
func (c *Client) GetRaids(ctx context.Context, e Expansion) (*Raids, error) {
	params := url.Values{"expansion_id": {strconv.Itoa(int(e))}}
	return getJSON[Raids](c, ctx, "/raiding/static-data", params)
}

// GetRaidRankings retrieves a list of raid rankings from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the RaidRankings struct
// Takes a RaidRankingsQuery struct as a parameter, in addition to context.Context
func (c *Client) GetRaidRankings(ctx context.Context, q *RaidRankingsQuery) (*RaidRankings, error) {
	return getJSONFor[RaidRankings](c, ctx, "/raiding/raid-rankings", q)
}

// GetGuildBossKill returns a guild's first kill of a given boss
// Takes a context.Context object to facilitate timeout, and a GuildBossKillQuery
// GuildBossKillQuery has only required fields for this request
// returns a BossKill object
func (c *Client) GetGuildBossKill(ctx context.Context, q *GuildBossKillQuery) (*BossKill, error) {
	return getJSONForMapped(c, ctx, "/guilds/boss-kill", q, mapBossKill)
}

// GetBossRankings retrieves the boss rankings for a given raid and boss
func (c *Client) GetBossRankings(ctx context.Context, q *BossRankingsQuery) (*BossRankings, error) {
	return getJSONFor[BossRankings](c, ctx, "/raiding/boss-rankings", q)
}

// GetHallOfFame retrieves the hall of fame for a given raid
func (c *Client) GetHallOfFame(ctx context.Context, q *RaidQuery) (*HallOfFame, error) {
	return getJSONFor[HallOfFame](c, ctx, "/raiding/hall-of-fame", q)
}

// GetRaidProgression retrieves the raid progression for a given raid
func (c *Client) GetRaidProgression(ctx context.Context, q *RaidQuery) (*RaidProgressionResponse, error) {
	return getJSONFor[RaidProgressionResponse](c, ctx, "/raiding/progression", q)
}

// GetMythicPlusAffixes retrieves the current mythic plus affixes for a region
func (c *Client) GetMythicPlusAffixes(ctx context.Context, q *AffixesQuery) (*MythicPlusAffixes, error) {
	return getJSONFor[MythicPlusAffixes](c, ctx, "/mythic-plus/affixes", q)
}

// GetMythicPlusStaticData retrieves the mythic plus seasons and dungeons for an expansion
func (c *Client) GetMythicPlusStaticData(ctx context.Context, e Expansion) (*MythicPlusStaticData, error) {
	params := url.Values{"expansion_id": {strconv.Itoa(int(e))}}
	return getJSON[MythicPlusStaticData](c, ctx, "/mythic-plus/static-data", params)
}

// GetMythicPlusRuns retrieves the mythic plus runs leaderboard.
// All query fields are optional; the API applies its own defaults when omitted
func (c *Client) GetMythicPlusRuns(ctx context.Context, q *MythicPlusRunsQuery) (*MythicPlusRuns, error) {
	return getJSONFor[MythicPlusRuns](c, ctx, "/mythic-plus/runs", q)
}

// GetMythicPlusRunDetails retrieves the details of a single mythic plus run
func (c *Client) GetMythicPlusRunDetails(ctx context.Context, q *RunDetailsQuery) (*MythicPlusRunEntry, error) {
	return getJSONFor[MythicPlusRunEntry](c, ctx, "/mythic-plus/run-details", q)
}

// GetMythicPlusScoreTiers retrieves the score/color breakpoints for a season.
// Season is optional; the API defaults to the current season when omitted
func (c *Client) GetMythicPlusScoreTiers(ctx context.Context, q *ScoreTiersQuery) ([]ScoreTier, error) {
	tiers, err := getJSONFor[[]ScoreTier](c, ctx, "/mythic-plus/score-tiers", q)
	if err != nil {
		return nil, err
	}
	return *tiers, nil
}

// GetMythicPlusSeasonCutoffs retrieves the score cutoffs for a season and region
func (c *Client) GetMythicPlusSeasonCutoffs(ctx context.Context, q *SeasonCutoffsQuery) (*SeasonCutoffs, error) {
	return getJSONFor[SeasonCutoffs](c, ctx, "/mythic-plus/season-cutoffs", q)
}

// GetMythicPlusLeaderboardCapacity retrieves the per-realm leaderboard capacity for a region
func (c *Client) GetMythicPlusLeaderboardCapacity(ctx context.Context, q *LeaderboardCapacityQuery) (*LeaderboardCapacity, error) {
	return getJSONFor[LeaderboardCapacity](c, ctx, "/mythic-plus/leaderboard-capacity", q)
}

// GetPeriods retrieves the current weekly reset periods for each region
func (c *Client) GetPeriods(ctx context.Context) (*Periods, error) {
	return getJSON[Periods](c, ctx, "/periods", url.Values{})
}
